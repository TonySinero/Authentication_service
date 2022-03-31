package repository

import (
	"database/sql"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/model"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/pkg/logging"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go

type AppUser interface {
	GetUserByID(id int) (*model.ResponseUser, error)
	GetUserAll(page int, limit int) ([]model.ResponseUser, int, error)
	GetUserByRoleFilter(page int, limit int, filters *model.RequestFilters) ([]model.ResponseUser, int, error)
	GetUserByDataFilter(page int, limit int, filters *model.RequestFilters) ([]model.ResponseUser, int, error)
	CreateStaff(User *model.CreateStaff) (int, error)
	CreateCustomer(User *model.CreateCustomer) (int, error)
	UpdateUser(User *model.UpdateUser) error
	DeleteUserByID(id int) (int, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUserPasswordByID(id int) (string, error)
	CheckEmail(email string) error
	RestorePassword(restore *model.RestorePassword) error
}

type Repository struct {
	AppUser
}

func NewRepository(db *sql.DB, logger logging.Logger) *Repository {
	return &Repository{
		AppUser: NewUserPostgres(db, logger),
	}
}
