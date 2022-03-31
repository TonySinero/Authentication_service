package repository

import (
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/model"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/pkg/logging"
	"testing"
	"time"
)

var logger = logging.GetLogger()

func TestRepository_GetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name          string
		mock          func(id int)
		id            int
		expectedUser  *model.ResponseUser
		expectedError bool
	}{
		{
			name: "OK",
			mock: func(id int) {
				rows := sqlmock.NewRows([]string{"id", "email", "role", "created_at"}).
					AddRow(1, "test@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)})
				mock.ExpectQuery("SELECT id, email, role, created_at FROM users WHERE id = (.+)").
					WithArgs(id).WillReturnRows(rows)
			},
			id: 1,
			expectedUser: &model.ResponseUser{
				ID:        1,
				Email:     "test@yandex.ru",
				Role:      "Courier",
				CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
			},
			expectedError: false,
		},
		{
			name: "Not found",
			mock: func(id int) {
				rows := sqlmock.NewRows([]string{"id", "email", "role", "created_at"})
				mock.ExpectQuery("SELECT id, email, role, created_at FROM users WHERE id = (.+)").
					WithArgs(id).WillReturnRows(rows)

			},
			id:            1,
			expectedUser:  nil,
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.id)
			got, err := r.GetUserByID(tt.id)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetUserPasswordByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name             string
		mock             func(id int)
		id               int
		expectedPassword string
		expectedError    bool
	}{
		{
			name: "OK",
			mock: func(id int) {
				rows := sqlmock.NewRows([]string{"password"}).
					AddRow("TestPassword")
				mock.ExpectQuery("SELECT password FROM users WHERE id = (.+)").
					WithArgs(id).WillReturnRows(rows)
			},
			id:               1,
			expectedPassword: "TestPassword",
			expectedError:    false,
		},
		{
			name: "Not found",
			mock: func(id int) {
				rows := sqlmock.NewRows([]string{"password"})
				mock.ExpectQuery("SELECT password FROM users WHERE id = (.+)").
					WithArgs(id).WillReturnRows(rows)
			},
			id:               1,
			expectedPassword: "",
			expectedError:    true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.id)
			got, err := r.GetUserPasswordByID(tt.id)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPassword, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetUserAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name          string
		inputPage     int
		inputLimit    int
		mock          func(page, limit int)
		expectedUser  []model.ResponseUser
		expectedPages int
		expectedError bool
	}{
		{
			name:       "Zero page and limit",
			inputPage:  0,
			inputLimit: 0,
			mock: func(page, limit int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id", "email", "role", "created_at"}).
					AddRow(1, "test@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)}).
					AddRow(2, "test1@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)}).
					AddRow(3, "test2@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)})
				mock.ExpectQuery("SELECT id, email, role, created_at FROM users").WillReturnRows(rows)
				mock.ExpectCommit()
			},

			expectedUser: []model.ResponseUser{
				{
					ID:        1,
					Email:     "test@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
				{
					ID:        2,
					Email:     "test1@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
				{
					ID:        3,
					Email:     "test2@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectedError: false,
		},
		{
			name:       "OK",
			inputPage:  1,
			inputLimit: 10,
			mock: func(page, limit int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id", "email", "role", "created_at"}).
					AddRow(1, "test@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)}).
					AddRow(2, "test1@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)}).
					AddRow(3, "test2@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)})

				mock.ExpectQuery("SELECT id, email, role, created_at FROM users WHERE deleted = false ORDER BY id LIMIT (.+) OFFSET (.+)").WithArgs(limit, (page-1)*limit).WillReturnRows(rows)
				rows2 := sqlmock.NewRows([]string{"pages"}).
					AddRow(1)
				mock.ExpectQuery("SELECT CEILING").WillReturnRows(rows2)
				mock.ExpectCommit()
			},

			expectedUser: []model.ResponseUser{
				{
					ID:        1,
					Email:     "test@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
				{
					ID:        2,
					Email:     "test1@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
				{
					ID:        3,
					Email:     "test2@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectedError: false,
		},
		{
			name:       "db error",
			inputPage:  1,
			inputLimit: 10,
			mock: func(page, limit int) {
				mock.ExpectBegin().WillReturnError(errors.New("some error"))
			},
			expectedUser:  nil,
			expectedError: true,
		},
		{
			name:       "db error2",
			inputPage:  1,
			inputLimit: 10,
			mock: func(page, limit int) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id, email, role, created_at FROM users WHERE deleted = false ORDER BY id LIMIT (.+) OFFSET (.+)").WithArgs(limit, (page-1)*limit).WillReturnError(errors.New("some error"))
			},
			expectedUser:  nil,
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputPage, tt.inputLimit)
			got, _, err := r.GetUserAll(tt.inputPage, tt.inputLimit)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetUserByRoleFilter(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name          string
		inputPage     int
		inputLimit    int
		inputFilter   *model.RequestFilters
		mock          func(page, limit int, filter *model.RequestFilters)
		expectedUser  []model.ResponseUser
		expectedPages int
		expectedError bool
	}{
		{
			name:       "OK Zero page and limit",
			inputPage:  0,
			inputLimit: 0,
			inputFilter: &model.RequestFilters{
				ShowDeleted: false,
				FilterData:  false,
				StartTime:   model.MyTime{},
				EndTime:     model.MyTime{},
				Role:        "Courier",
			},
			mock: func(page, limit int, filter *model.RequestFilters) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id", "email", "role", "created_at"}).
					AddRow(1, "test@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)}).
					AddRow(2, "test1@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)}).
					AddRow(3, "test2@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)})
				mock.ExpectQuery("SELECT id, email, role, created_at FROM users WHERE").WithArgs(filter.Role).WillReturnRows(rows)
				mock.ExpectCommit()
			},

			expectedUser: []model.ResponseUser{
				{
					ID:        1,
					Email:     "test@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
				{
					ID:        2,
					Email:     "test1@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
				{
					ID:        3,
					Email:     "test2@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectedError: false,
		},
		{
			name:       "OK With page and limit",
			inputPage:  1,
			inputLimit: 10,
			inputFilter: &model.RequestFilters{
				ShowDeleted: false,
				FilterData:  false,
				StartTime:   model.MyTime{},
				EndTime:     model.MyTime{},
				Role:        "Courier",
			},
			mock: func(page, limit int, filter *model.RequestFilters) {
				mock.ExpectBegin()
				rowsForPages := sqlmock.NewRows([]string{"pages"}).AddRow("1")
				mock.ExpectQuery("SELECT CEILING").WithArgs(limit, filter.Role).WillReturnRows(rowsForPages)
				rows := sqlmock.NewRows([]string{"id", "email", "role", "created_at"}).
					AddRow(1, "test@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)}).
					AddRow(2, "test1@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)}).
					AddRow(3, "test2@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)})
				mock.ExpectQuery("SELECT id, email, role, created_at FROM users WHERE").WithArgs(filter.Role, limit, (page-1)*limit).WillReturnRows(rows)
				mock.ExpectCommit()
			},

			expectedUser: []model.ResponseUser{
				{
					ID:        1,
					Email:     "test@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
				{
					ID:        2,
					Email:     "test1@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
				{
					ID:        3,
					Email:     "test2@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectedError: false,
		},
		{
			name:       "OK showDeleted = true Zero page and limit",
			inputPage:  0,
			inputLimit: 0,
			inputFilter: &model.RequestFilters{
				ShowDeleted: true,
				FilterData:  false,
				StartTime:   model.MyTime{},
				EndTime:     model.MyTime{},
				Role:        "Courier",
			},
			mock: func(page, limit int, filter *model.RequestFilters) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id", "email", "role", "created_at"}).
					AddRow(1, "test@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)}).
					AddRow(2, "test1@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)}).
					AddRow(3, "test2@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)})
				mock.ExpectQuery("SELECT id, email, role, created_at FROM users WHERE").WithArgs(filter.Role).WillReturnRows(rows)
				mock.ExpectCommit()
			},

			expectedUser: []model.ResponseUser{
				{
					ID:        1,
					Email:     "test@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
				{
					ID:        2,
					Email:     "test1@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
				{
					ID:        3,
					Email:     "test2@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectedError: false,
		},
		{
			name:       "db error",
			inputPage:  1,
			inputLimit: 10,
			mock: func(page, limit int, filter *model.RequestFilters) {
				mock.ExpectBegin().WillReturnError(errors.New("some error"))
			},
			expectedUser:  nil,
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputPage, tt.inputLimit, tt.inputFilter)
			got, _, err := r.GetUserByRoleFilter(tt.inputPage, tt.inputLimit, tt.inputFilter)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetUserByDataFilter(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name          string
		inputPage     int
		inputLimit    int
		inputFilter   *model.RequestFilters
		mock          func(page, limit int, filter *model.RequestFilters)
		expectedUser  []model.ResponseUser
		expectedPages int
		expectedError bool
	}{
		{
			name:       "OK Zero page and limit",
			inputPage:  0,
			inputLimit: 0,
			inputFilter: &model.RequestFilters{
				ShowDeleted: false,
				FilterData:  true,
				StartTime:   model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				EndTime:     model.MyTime{},
				Role:        "",
			},
			mock: func(page, limit int, filter *model.RequestFilters) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id", "email", "role", "created_at"}).
					AddRow(1, "test@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)}).
					AddRow(2, "test1@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)}).
					AddRow(3, "test2@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)})
				mock.ExpectQuery("SELECT id, email, role, created_at FROM users WHERE").WithArgs(filter.StartTime, filter.EndTime).WillReturnRows(rows)
				mock.ExpectCommit()
			},

			expectedUser: []model.ResponseUser{
				{
					ID:        1,
					Email:     "test@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
				{
					ID:        2,
					Email:     "test1@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
				{
					ID:        3,
					Email:     "test2@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectedError: false,
		},
		{
			name:       "OK With page and limit",
			inputPage:  1,
			inputLimit: 10,
			inputFilter: &model.RequestFilters{
				ShowDeleted: false,
				FilterData:  true,
				StartTime:   model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				EndTime:     model.MyTime{},
				Role:        "",
			},
			mock: func(page, limit int, filter *model.RequestFilters) {
				mock.ExpectBegin()
				rowsForPages := sqlmock.NewRows([]string{"pages"}).AddRow("1")
				mock.ExpectQuery("SELECT CEILING").WithArgs(limit, filter.StartTime, filter.EndTime).WillReturnRows(rowsForPages)
				rows := sqlmock.NewRows([]string{"id", "email", "role", "created_at"}).
					AddRow(1, "test@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)}).
					AddRow(2, "test1@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)}).
					AddRow(3, "test2@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)})
				mock.ExpectQuery("SELECT id, email, role, created_at FROM users WHERE").WithArgs(filter.StartTime, filter.EndTime, limit, (page-1)*limit).WillReturnRows(rows)
				mock.ExpectCommit()
			},

			expectedUser: []model.ResponseUser{
				{
					ID:        1,
					Email:     "test@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
				{
					ID:        2,
					Email:     "test1@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
				{
					ID:        3,
					Email:     "test2@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectedError: false,
		},
		{
			name:       "OK showDeleted = true Zero page and limit",
			inputPage:  0,
			inputLimit: 0,
			inputFilter: &model.RequestFilters{
				ShowDeleted: true,
				FilterData:  true,
				StartTime:   model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				EndTime:     model.MyTime{},
				Role:        "",
			},
			mock: func(page, limit int, filter *model.RequestFilters) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id", "email", "role", "created_at"}).
					AddRow(1, "test@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)}).
					AddRow(2, "test1@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)}).
					AddRow(3, "test2@yandex.ru", "Courier", model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)})
				mock.ExpectQuery("SELECT id, email, role, created_at FROM users WHERE").WithArgs(filter.StartTime, filter.EndTime).WillReturnRows(rows)
				mock.ExpectCommit()
			},

			expectedUser: []model.ResponseUser{
				{
					ID:        1,
					Email:     "test@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
				{
					ID:        2,
					Email:     "test1@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
				{
					ID:        3,
					Email:     "test2@yandex.ru",
					Role:      "Courier",
					CreatedAt: model.MyTime{Time: time.Date(2022, 03, 11, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectedError: false,
		},
		{
			name:       "db error",
			inputPage:  1,
			inputLimit: 10,
			mock: func(page, limit int, filter *model.RequestFilters) {
				mock.ExpectBegin().WillReturnError(errors.New("some error"))
			},
			expectedUser:  nil,
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputPage, tt.inputLimit, tt.inputFilter)
			got, _, err := r.GetUserByDataFilter(tt.inputPage, tt.inputLimit, tt.inputFilter)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name          string
		mock          func(email string)
		email         string
		expectedUser  *model.User
		expectedError bool
	}{
		{
			name: "OK",
			mock: func(email string) {
				rows := sqlmock.NewRows([]string{"id", "email", "password", "role", "deleted"}).
					AddRow(1, "test@yandex.ru", "$2a$10$ooCmcWnLIubagB1MqM3UWOIpJTrq58tPQO6HVraj3yTKASiXBXHqy", "Courier", false)

				mock.ExpectQuery("SELECT id, email, password, role, deleted FROM users WHERE email = (.+)").
					WithArgs(email).WillReturnRows(rows)
			},
			email: "test@yandex.ru",
			expectedUser: &model.User{
				ID:       1,
				Email:    "test@yandex.ru",
				Password: "$2a$10$ooCmcWnLIubagB1MqM3UWOIpJTrq58tPQO6HVraj3yTKASiXBXHqy",
				Role:     "Courier",
				Deleted:  false,
			},
			expectedError: false,
		},
		{
			name: "Not found",
			mock: func(email string) {
				rows := sqlmock.NewRows([]string{"id", "email", "password", "role", "deleted"})

				mock.ExpectQuery("SELECT id, email, password, role, deleted FROM users WHERE email = (.+)").
					WithArgs(email).WillReturnRows(rows).WillReturnError(errors.New("some error"))

			},
			email:         "test@yandex.ru",
			expectedUser:  nil,
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.email)
			got, err := r.GetUserByEmail(tt.email)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_DeleteUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name           string
		mock           func(id int)
		id             int
		expectedUserId int
		expectedError  bool
	}{
		{
			name: "OK",
			mock: func(id int) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1)
				mock.ExpectQuery("UPDATE users SET deleted = true WHERE id=(.+) RETURNING id").
					WithArgs(id).WillReturnRows(rows)
			},
			id:             1,
			expectedUserId: 1,
			expectedError:  false,
		},
		{
			name: "Not found",
			mock: func(id int) {
				rows := sqlmock.NewRows([]string{"id"})

				mock.ExpectQuery("UPDATE users SET deleted = true WHERE id=(.+) RETURNING id").
					WithArgs(id).WillReturnRows(rows)

			},
			id:             1,
			expectedUserId: 0,
			expectedError:  true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.id)
			got, err := r.DeleteUserByID(tt.id)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUserId, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_CreateStaff(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name           string
		mock           func(user *model.CreateStaff)
		InputUser      *model.CreateStaff
		expectedUserId int
		expectedError  bool
	}{
		{
			name: "OK",
			mock: func(user *model.CreateStaff) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1)
				mock.ExpectQuery("INSERT INTO users").WithArgs(user.Email, user.Password, user.Role, time.Now().Format(model.Layout), false).
					WillReturnRows(rows)
			},
			InputUser: &model.CreateStaff{
				Email:    "test@yandex.ru",
				Password: "$2a$10$EpAGhm0HGkxBiPyBAB7xzuyEbZlZCjvSdcJTjamaJyxZRir1vaMmW",
				Role:     "Courier",
			},
			expectedUserId: 1,
			expectedError:  false,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.InputUser)
			got, err := r.CreateStaff(tt.InputUser)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUserId, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_CreateCustomer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name           string
		mock           func(user *model.CreateCustomer)
		InputUser      *model.CreateCustomer
		expectedUserId int
		expectedError  bool
	}{
		{
			name: "OK",
			mock: func(user *model.CreateCustomer) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1)
				mock.ExpectQuery("INSERT INTO users").WithArgs(user.Email, user.Password, "Authorized Customer", time.Now().Format(model.Layout), false).
					WillReturnRows(rows)
			},
			InputUser: &model.CreateCustomer{
				Email:    "test@yandex.ru",
				Password: "$2a$10$EpAGhm0HGkxBiPyBAB7xzuyEbZlZCjvSdcJTjamaJyxZRir1vaMmW",
			},
			expectedUserId: 1,
			expectedError:  false,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.InputUser)
			got, err := r.CreateCustomer(tt.InputUser)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUserId, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_UpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name          string
		mock          func(user *model.UpdateUser)
		InputUser     *model.UpdateUser
		expectedError bool
	}{
		{
			name: "OK",
			mock: func(user *model.UpdateUser) {
				result := driver.ResultNoRows
				mock.ExpectExec("UPDATE users").
					WithArgs(user.NewPassword, user.Email).WillReturnResult(result)
			},
			InputUser: &model.UpdateUser{
				Email:       "test@yandex.ru",
				OldPassword: "$2a$10$EpAGhm0HGkxBiPyBAB7xzuyEbZlZCjvSdcJTjamaJyxZRir1vaMmW",
				NewPassword: "$2a$10$EpAGhm0HGkxBiPyBAB7xzuyEbZlZCjvSdcJTjamaJyxZRir1vaMmW",
			},
			expectedError: false,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.InputUser)
			err := r.UpdateUser(tt.InputUser)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
