package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/GRPC/grpcClient"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/model"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/pkg/logging"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/repository"
	mock_repository "stlab.itechart-group.com/go/food_delivery/authentication_service/repository/mocks"
	"testing"
	"time"
)

func TestService_GetUser(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockAppUser, id int)
	testTable := []struct {
		name          string
		inputId       int
		mockBehavior  mockBehavior
		expectedUser  *model.ResponseUser
		expectedError error
	}{
		{
			name:    "OK",
			inputId: 1,
			mockBehavior: func(s *mock_repository.MockAppUser, id int) {
				s.EXPECT().GetUserByID(id).Return(&model.ResponseUser{
					ID:        1,
					Email:     "test@yandex.ru",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
					Role:      "Courier",
				}, nil)
			},
			expectedUser: &model.ResponseUser{
				ID:        1,
				Email:     "test@yandex.ru",
				CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				Role:      "Courier",
			},
			expectedError: nil,
		},
		{
			name:    "Repository failure",
			inputId: 1,
			mockBehavior: func(s *mock_repository.MockAppUser, id int) {
				s.EXPECT().GetUserByID(id).Return(nil, errors.New("repository failure"))
			},
			expectedUser:  nil,
			expectedError: errors.New("repository failure"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			auth := mock_repository.NewMockAppUser(c)
			testCase.mockBehavior(auth, testCase.inputId)
			logger := logging.GetLogger()
			repo := &repository.Repository{AppUser: auth}
			grpcCli := grpcClient.NewGRPCClient("159.223.1.135")
			service := NewService(repo, grpcCli, logger)
			user, err := service.GetUser(testCase.inputId)
			//Assert
			assert.Equal(t, testCase.expectedUser, user)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestService_GetUsers(t *testing.T) {
	type mockBehaviorGetAllUsers func(s *mock_repository.MockAppUser, page int, limit int)
	type mockBehaviorFilterByRole func(s *mock_repository.MockAppUser, page int, limit int, filter *model.RequestFilters)
	type mockBehaviorFilterByData func(s *mock_repository.MockAppUser, page int, limit int, filter *model.RequestFilters)
	testTable := []struct {
		name                     string
		inputPage                int
		inputLimit               int
		inputFilter              *model.RequestFilters
		mockBehaviorGetAllUsers  mockBehaviorGetAllUsers
		mockBehaviorFilterByRole mockBehaviorFilterByRole
		mockBehaviorFilterByData mockBehaviorFilterByData
		expectedUsers            []model.ResponseUser
		expectedError            error
	}{
		{
			name:       "OK without filter",
			inputPage:  1,
			inputLimit: 10,
			inputFilter: &model.RequestFilters{
				ShowDeleted: false,
				FilterData:  false,
				StartTime:   model.MyTime{},
				EndTime:     model.MyTime{},
				Role:        "",
			},
			mockBehaviorGetAllUsers: func(s *mock_repository.MockAppUser, page int, limit int) {
				s.EXPECT().GetUserAll(page, limit).Return([]model.ResponseUser{
					{ID: 1,
						Email:     "test@yande.ru",
						CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
						Role:      "Courier",
					}, {ID: 2,
						Email:     "test2@yande.ru",
						CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
						Role:      "Courier",
					},
				}, 1, nil)
			},
			mockBehaviorFilterByRole: func(s *mock_repository.MockAppUser, page int, limit int, filter *model.RequestFilters) {},
			mockBehaviorFilterByData: func(s *mock_repository.MockAppUser, page int, limit int, filter *model.RequestFilters) {},
			expectedUsers: []model.ResponseUser{
				{ID: 1,
					Email:     "test@yande.ru",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
					Role:      "Courier",
				}, {ID: 2,
					Email:     "test2@yande.ru",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
					Role:      "Courier",
				},
			},
			expectedError: nil,
		},
		{
			name:       "Repository failure",
			inputPage:  1,
			inputLimit: 10,
			inputFilter: &model.RequestFilters{
				ShowDeleted: false,
				FilterData:  false,
				StartTime:   model.MyTime{},
				EndTime:     model.MyTime{},
				Role:        "",
			},
			mockBehaviorGetAllUsers: func(s *mock_repository.MockAppUser, page int, limit int) {
				s.EXPECT().GetUserAll(page, limit).Return(nil, 0, errors.New("repository failure"))
			},
			mockBehaviorFilterByRole: func(s *mock_repository.MockAppUser, page int, limit int, filter *model.RequestFilters) {},
			mockBehaviorFilterByData: func(s *mock_repository.MockAppUser, page int, limit int, filter *model.RequestFilters) {},
			expectedUsers:            nil,
			expectedError:            errors.New("repository failure"),
		},
		{
			name:       "OK with role filter",
			inputPage:  1,
			inputLimit: 10,
			inputFilter: &model.RequestFilters{
				ShowDeleted: false,
				FilterData:  false,
				StartTime:   model.MyTime{},
				EndTime:     model.MyTime{},
				Role:        "Courier",
			},
			mockBehaviorGetAllUsers: func(s *mock_repository.MockAppUser, page int, limit int) {},
			mockBehaviorFilterByRole: func(s *mock_repository.MockAppUser, page int, limit int, filter *model.RequestFilters) {
				s.EXPECT().GetUserByRoleFilter(page, limit, filter).Return([]model.ResponseUser{
					{ID: 1,
						Email:     "test@yande.ru",
						CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
						Role:      "Courier",
					}, {ID: 2,
						Email:     "test2@yande.ru",
						CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
						Role:      "Courier",
					},
				}, 1, nil)
			},
			mockBehaviorFilterByData: func(s *mock_repository.MockAppUser, page int, limit int, filter *model.RequestFilters) {},
			expectedUsers: []model.ResponseUser{
				{ID: 1,
					Email:     "test@yande.ru",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
					Role:      "Courier",
				}, {ID: 2,
					Email:     "test2@yande.ru",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
					Role:      "Courier",
				},
			},
			expectedError: nil,
		},
		{
			name:       "Repository failure (role filter)",
			inputPage:  1,
			inputLimit: 10,
			inputFilter: &model.RequestFilters{
				ShowDeleted: false,
				FilterData:  false,
				StartTime:   model.MyTime{},
				EndTime:     model.MyTime{},
				Role:        "Courier",
			},
			mockBehaviorGetAllUsers: func(s *mock_repository.MockAppUser, page int, limit int) {},
			mockBehaviorFilterByRole: func(s *mock_repository.MockAppUser, page int, limit int, filter *model.RequestFilters) {
				s.EXPECT().GetUserByRoleFilter(page, limit, filter).Return(nil, 0, errors.New("repository failure"))
			},
			mockBehaviorFilterByData: func(s *mock_repository.MockAppUser, page int, limit int, filter *model.RequestFilters) {},
			expectedUsers:            nil,
			expectedError:            errors.New("repository failure"),
		},
		{
			name:       "OK with data filter",
			inputPage:  1,
			inputLimit: 10,
			inputFilter: &model.RequestFilters{
				ShowDeleted: false,
				FilterData:  true,
				StartTime:   model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				EndTime:     model.MyTime{},
				Role:        "",
			},
			mockBehaviorGetAllUsers:  func(s *mock_repository.MockAppUser, page int, limit int) {},
			mockBehaviorFilterByRole: func(s *mock_repository.MockAppUser, page int, limit int, filter *model.RequestFilters) {},
			mockBehaviorFilterByData: func(s *mock_repository.MockAppUser, page int, limit int, filter *model.RequestFilters) {
				s.EXPECT().GetUserByDataFilter(page, limit, filter).Return([]model.ResponseUser{
					{ID: 1,
						Email:     "test@yande.ru",
						CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
						Role:      "Courier",
					}, {ID: 2,
						Email:     "test2@yande.ru",
						CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
						Role:      "Courier",
					},
				}, 1, nil)
			},
			expectedUsers: []model.ResponseUser{
				{ID: 1,
					Email:     "test@yande.ru",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
					Role:      "Courier",
				}, {ID: 2,
					Email:     "test2@yande.ru",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
					Role:      "Courier",
				},
			},
			expectedError: nil,
		},
		{
			name:       "Repository failure (data filter)",
			inputPage:  1,
			inputLimit: 10,
			inputFilter: &model.RequestFilters{
				ShowDeleted: false,
				FilterData:  true,
				StartTime:   model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				EndTime:     model.MyTime{},
				Role:        "",
			},
			mockBehaviorGetAllUsers:  func(s *mock_repository.MockAppUser, page int, limit int) {},
			mockBehaviorFilterByRole: func(s *mock_repository.MockAppUser, page int, limit int, filter *model.RequestFilters) {},
			mockBehaviorFilterByData: func(s *mock_repository.MockAppUser, page int, limit int, filter *model.RequestFilters) {
				s.EXPECT().GetUserByDataFilter(page, limit, filter).Return(nil, 0, errors.New("repository failure"))
			},
			expectedUsers: nil,
			expectedError: errors.New("repository failure"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			auth := mock_repository.NewMockAppUser(c)
			testCase.mockBehaviorGetAllUsers(auth, testCase.inputPage, testCase.inputLimit)
			testCase.mockBehaviorFilterByRole(auth, testCase.inputPage, testCase.inputLimit, testCase.inputFilter)
			testCase.mockBehaviorFilterByData(auth, testCase.inputPage, testCase.inputLimit, testCase.inputFilter)
			logger := logging.GetLogger()
			repo := &repository.Repository{AppUser: auth}

			grpcCli := grpcClient.NewGRPCClient("159.223.1.135")
			service := NewService(repo, grpcCli, logger)
			users, _, err := service.GetUsers(testCase.inputPage, testCase.inputLimit, testCase.inputFilter)
			//Assert
			assert.Equal(t, testCase.expectedUsers, users)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestService_UpdateUser(t *testing.T) {
	type mockBehaviorUpdate func(s *mock_repository.MockAppUser, user *model.UpdateUser)
	type mockBehaviorGet func(s *mock_repository.MockAppUser, user *model.UpdateUser)
	testTable := []struct {
		name               string
		inputUser          *model.UpdateUser
		mockBehaviorUpdate mockBehaviorUpdate
		mockBehaviorGet    mockBehaviorGet
		expectedError      error
	}{
		{
			name: "OK",
			inputUser: &model.UpdateUser{
				Email:       "test@yandex.ru",
				OldPassword: "HGYKnu!98Tg",
				NewPassword: "HYKnu!98Tg",
			},
			mockBehaviorUpdate: func(s *mock_repository.MockAppUser, user *model.UpdateUser) {
				s.EXPECT().UpdateUser(user).Return(nil)
			},
			mockBehaviorGet: func(s *mock_repository.MockAppUser, user *model.UpdateUser) {
				s.EXPECT().GetUserByEmail(user.Email).Return(&model.User{
					ID:       1,
					Email:    "test@yandex.ru",
					Password: "$2a$10$ooCmcWnLIubagB1MqM3UWOIpJTrq58tPQO6HVraj3yTKASiXBXHqy",
					Role:     "Courier",
				}, nil)
			},
			expectedError: nil,
		},
		{
			name: "Error while getting user",
			inputUser: &model.UpdateUser{
				Email:       "test@yandex.ru",
				OldPassword: "HGYKnu!98Tg",
				NewPassword: "HYKnu!98Tg",
			},
			mockBehaviorUpdate: func(s *mock_repository.MockAppUser, user *model.UpdateUser) {},
			mockBehaviorGet: func(s *mock_repository.MockAppUser, user *model.UpdateUser) {
				s.EXPECT().GetUserByEmail(user.Email).Return(nil, errors.New("error while getting user"))
			},
			expectedError: errors.New("error while getting user"),
		},
		{
			name: "Error while updating user",
			inputUser: &model.UpdateUser{
				Email:       "test@yandex.ru",
				OldPassword: "HGYKnu!98Tg",
				NewPassword: "HYKnu!98Tg",
			},
			mockBehaviorUpdate: func(s *mock_repository.MockAppUser, user *model.UpdateUser) {
				s.EXPECT().UpdateUser(user).Return(errors.New("error while getting user"))
			},
			mockBehaviorGet: func(s *mock_repository.MockAppUser, user *model.UpdateUser) {
				s.EXPECT().GetUserByEmail(user.Email).Return(&model.User{
					ID:       1,
					Email:    "test@yandex.ru",
					Password: "$2a$10$ooCmcWnLIubagB1MqM3UWOIpJTrq58tPQO6HVraj3yTKASiXBXHqy",
					Role:     "Courier",
				}, nil)
			},

			expectedError: errors.New("error while getting user"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			auth := mock_repository.NewMockAppUser(c)
			testCase.mockBehaviorUpdate(auth, testCase.inputUser)
			testCase.mockBehaviorGet(auth, testCase.inputUser)
			logger := logging.GetLogger()
			repo := &repository.Repository{AppUser: auth}
			grpcCli := grpcClient.NewGRPCClient("159.223.1.135")
			service := NewService(repo, grpcCli, logger)
			err := service.UpdateUser(testCase.inputUser)
			//Assert

			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestService_DeleteUser(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockAppUser, id int)
	testTable := []struct {
		name           string
		inputId        int
		mockBehavior   mockBehavior
		expectedUserId int
		expectedError  error
	}{
		{
			name:    "OK",
			inputId: 1,
			mockBehavior: func(s *mock_repository.MockAppUser, id int) {
				s.EXPECT().DeleteUserByID(id).Return(1, nil)
			},
			expectedUserId: 1,
			expectedError:  nil,
		},
		{
			name:    "Repository failure",
			inputId: 1,
			mockBehavior: func(s *mock_repository.MockAppUser, id int) {
				s.EXPECT().DeleteUserByID(id).Return(0, errors.New("repository failure"))
			},
			expectedUserId: 0,
			expectedError:  errors.New("repository failure"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			auth := mock_repository.NewMockAppUser(c)
			testCase.mockBehavior(auth, testCase.inputId)
			logger := logging.GetLogger()
			repo := &repository.Repository{AppUser: auth}
			grpcCli := grpcClient.NewGRPCClient("159.223.1.135")
			service := NewService(repo, grpcCli, logger)
			id, err := service.DeleteUserByID(testCase.inputId)
			//Assert
			assert.Equal(t, testCase.expectedUserId, id)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestService_RestorePassword(t *testing.T) {
	type mockBehaviorCheckEmail func(s *mock_repository.MockAppUser, email string)
	type mockBehavior func(s *mock_repository.MockAppUser, restore *model.RestorePassword)
	testTable := []struct {
		name                   string
		input                  *model.RestorePassword
		mockBehaviorCheckEmail mockBehaviorCheckEmail
		mockBehavior           mockBehavior
		expectedError          error
	}{
		{
			name: "OK",
			input: &model.RestorePassword{
				Email: "test@yandex.ru",
			},
			mockBehaviorCheckEmail: func(s *mock_repository.MockAppUser, email string) {
				s.EXPECT().CheckEmail(email).Return(nil)
			},
			mockBehavior: func(s *mock_repository.MockAppUser, restore *model.RestorePassword) {
				s.EXPECT().RestorePassword(restore).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "user does not exist",
			input: &model.RestorePassword{
				Email: "test@yandex.ru",
			},
			mockBehaviorCheckEmail: func(s *mock_repository.MockAppUser, email string) {
				s.EXPECT().CheckEmail(email).Return(errors.New("user with this email does not exist"))
			},
			mockBehavior:  func(s *mock_repository.MockAppUser, restore *model.RestorePassword) {},
			expectedError: errors.New("user with this email does not exist"),
		},
		{
			name: "Error while updating user",
			input: &model.RestorePassword{
				Email: "test@yandex.ru",
			},
			mockBehaviorCheckEmail: func(s *mock_repository.MockAppUser, email string) {
				s.EXPECT().CheckEmail(email).Return(nil)
			},
			mockBehavior: func(s *mock_repository.MockAppUser, restore *model.RestorePassword) {
				s.EXPECT().RestorePassword(restore).Return(errors.New("error while updating password"))
			},
			expectedError: errors.New("error while updating password"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			auth := mock_repository.NewMockAppUser(c)
			testCase.mockBehaviorCheckEmail(auth, testCase.input.Email)
			testCase.mockBehavior(auth, testCase.input)
			logger := logging.GetLogger()
			repo := &repository.Repository{AppUser: auth}
			grpcCli := grpcClient.NewGRPCClient("159.223.1.135")
			service := NewService(repo, grpcCli, logger)
			err := service.RestorePassword(testCase.input)
			//Assert

			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
