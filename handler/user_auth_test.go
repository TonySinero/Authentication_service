package handler

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	authProto "stlab.itechart-group.com/go/food_delivery/authentication_service/GRPC"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/model"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/pkg/logging"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/service"
	mock_service "stlab.itechart-group.com/go/food_delivery/authentication_service/service/mocks"
	"testing"
)

func TestHandler_authUser(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAppUser, user model.AuthUser)
	testTable := []struct {
		name                string
		inputBody           string
		inputUser           model.AuthUser
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email":"test@yandex.ru", "password":"HGYKnu!98Tg"}`,
			inputUser: model.AuthUser{
				Email:    "test@yandex.ru",
				Password: "HGYKnu!98Tg",
			},
			mockBehavior: func(s *mock_service.MockAppUser, user model.AuthUser) {
				s.EXPECT().AuthUser(user.Email, user.Password).Return(&authProto.GeneratedTokens{
					AccessToken:  "qwerty",
					RefreshToken: "qwerty",
				}, 1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"accessToken":"qwerty","refreshToken":"qwerty"}`,
		},
		{
			name:                "Empty fields",
			inputBody:           `{"email":"test@yandex.ru"}`,
			inputUser:           model.AuthUser{},
			mockBehavior:        func(s *mock_service.MockAppUser, user model.AuthUser) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"Invalid input body"}`,
		},
		{
			name:      "invalid values in email field",
			inputBody: `{"email":"testyandex.ru", "password":"HGYKnu!98Tg"}`,
			inputUser: model.AuthUser{
				Email:    "testyandex.ru",
				Password: "HGYKnu!98Tg",
			},
			mockBehavior:        func(s *mock_service.MockAppUser, user model.AuthUser) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"Wrong email or password entered"}`,
		},
		{
			name:      "invalid values in password field",
			inputBody: `{"email":"test@yandex.ru", "password":"HGYKnu98Tg"}`,
			inputUser: model.AuthUser{
				Email:    "test@yandex.ru",
				Password: "HGYKnu98Tg",
			},
			mockBehavior:        func(s *mock_service.MockAppUser, user model.AuthUser) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"Wrong email or password entered"}`,
		},
		{
			name:      "invalid values in both fields",
			inputBody: `{"email":"testyandex.ru", "password":"HGYKnu98Tg"}`,
			inputUser: model.AuthUser{
				Email:    "testyandex.ru",
				Password: "HGYKnu98Tg",
			},
			mockBehavior:        func(s *mock_service.MockAppUser, user model.AuthUser) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"Wrong email or password entered"}`,
		},
		{
			name:      "invalid length of the password",
			inputBody: `{"email":"test@yandex.ru", "password":"Hnu!9T"}`,
			inputUser: model.AuthUser{
				Email:    "test@yandex.ru",
				Password: "Hnu!9T",
			},
			mockBehavior:        func(s *mock_service.MockAppUser, user model.AuthUser) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"Wrong email or password entered"}`,
		},
		{
			name:      "space in the password",
			inputBody: `{"email":"test@yandex.ru", "password":"Hn   u!9T"}`,
			inputUser: model.AuthUser{
				Email:    "test@yandex.ru",
				Password: "Hn   u!9T",
			},
			mockBehavior:        func(s *mock_service.MockAppUser, user model.AuthUser) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"Wrong email or password entered"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"email":"test@yandex.ru", "password":"HGYKnu!98Tg"}`,
			inputUser: model.AuthUser{
				Email:    "test@yandex.ru",
				Password: "HGYKnu!98Tg",
			},
			mockBehavior: func(s *mock_service.MockAppUser, user model.AuthUser) {
				s.EXPECT().AuthUser(user.Email, user.Password).Return(nil, 0, errors.New("service failure"))
			},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"Wrong email or password entered"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			auth := mock_service.NewMockAppUser(c)
			testCase.mockBehavior(auth, testCase.inputUser)
			logger := logging.GetLogger()
			services := &service.Service{AppUser: auth}
			handler := NewHandler(logger, services)

			//Init server
			r := gin.New()
			r.POST("/login", handler.authUser)

			//Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(testCase.inputBody))

			//Execute the request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}

}
