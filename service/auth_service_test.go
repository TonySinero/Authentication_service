package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	authProto "stlab.itechart-group.com/go/food_delivery/authentication_service/GRPC"
	mockAuthProto "stlab.itechart-group.com/go/food_delivery/authentication_service/GRPC/authProto"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/GRPC/grpcClient"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/model"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/pkg/logging"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/repository"
	mock_repository "stlab.itechart-group.com/go/food_delivery/authentication_service/repository/mocks"
	"testing"
)

func TestService_authUser(t *testing.T) {
	type mockBehaviorGetUser func(s *mock_repository.MockAppUser, email string)
	type mockBehaviorGetTokens func(s *mockAuthProto.MockAuthServer, user *authProto.User) (*authProto.GeneratedTokens, error)
	testTable := []struct {
		name                  string
		inputPassword         string
		inputEmail            string
		mockUser              *authProto.User
		mockBehaviorGetUser   mockBehaviorGetUser
		mockBehaviorGetTokens mockBehaviorGetTokens
		expectedId            int
		expectedError         error
	}{
		{
			name:          "OK",
			inputPassword: "HGYKnu!98Tg",
			inputEmail:    "test@yandex.ru",
			mockUser: &authProto.User{
				UserId: 1,
				Role:   "Superadmin",
			},
			mockBehaviorGetUser: func(s *mock_repository.MockAppUser, email string) {
				s.EXPECT().GetUserByEmail(email).Return(&model.User{
					ID:       1,
					Email:    "test@yandex.ru",
					Password: "$2a$10$ooCmcWnLIubagB1MqM3UWOIpJTrq58tPQO6HVraj3yTKASiXBXHqy",
					Deleted:  false,
				}, nil)
			},
			mockBehaviorGetTokens: func(s *mockAuthProto.MockAuthServer, user *authProto.User) (*authProto.GeneratedTokens, error) {
				_, _ = s.TokenGenerationByUserId(context.Background(), user)
				return &authProto.GeneratedTokens{AccessToken: "qwerty", RefreshToken: "qwerty"}, nil
			},
			expectedId:    1,
			expectedError: nil,
		},
		{
			name:          "Wrong password",
			inputPassword: "HGYKnu!9Tg",
			inputEmail:    "test@yandex.ru",
			mockBehaviorGetUser: func(s *mock_repository.MockAppUser, email string) {
				s.EXPECT().GetUserByEmail(email).Return(&model.User{
					ID:       1,
					Email:    "test@yandex.ru",
					Password: "$2a$10$ooCmcWnLIubagB1MqM3UWOIpJTrq58tPQO6HVraj3yTKASiXBXHqy",
					Deleted:  false,
				}, nil)
			},
			mockBehaviorGetTokens: func(s *mockAuthProto.MockAuthServer, user *authProto.User) (*authProto.GeneratedTokens, error) {
				return nil, nil
			},
			expectedError: errors.New("wrong email or password entered"),
		},
		{
			name:          "Repository error",
			inputPassword: "HGYKnu!98Tg",
			inputEmail:    "test@yandex.ru",
			mockBehaviorGetUser: func(s *mock_repository.MockAppUser, email string) {
				s.EXPECT().GetUserByEmail(email).Return(nil, errors.New("repository error"))
			},
			mockBehaviorGetTokens: func(s *mockAuthProto.MockAuthServer, user *authProto.User) (*authProto.GeneratedTokens, error) {
				return nil, nil
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			repo := mock_repository.NewMockAppUser(c)
			reposit := &repository.Repository{AppUser: repo}
			testCase.mockBehaviorGetUser(repo, testCase.inputEmail)
			mockProto := new(mockAuthProto.MockAuthServer)
			testCase.mockBehaviorGetTokens(mockProto, testCase.mockUser)
			logger := logging.GetLogger()
			grpcCli := grpcClient.NewGRPCClient("159.223.1.135")
			service := NewService(reposit, grpcCli, logger)
			_, id, err := service.AuthUser(testCase.inputEmail, testCase.inputPassword)
			//Assert
			assert.Equal(t, testCase.expectedId, id)
			assert.Equal(t, testCase.expectedError, err)
		})
	}

}
