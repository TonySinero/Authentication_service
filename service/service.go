package service

import (
	authProto "stlab.itechart-group.com/go/food_delivery/authentication_service/GRPC"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/GRPC/grpcClient"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/model"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/pkg/logging"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go

type AppUser interface {
	GetUser(id int) (*model.ResponseUser, error)
	GetUsers(page int, limit int, filters *model.RequestFilters) ([]model.ResponseUser, int, error)
	CreateCustomer(user *model.CreateCustomer) (*authProto.GeneratedTokens, int, error)
	CreateStaff(user *model.CreateStaff) (int, error)
	UpdateUser(user *model.UpdateUser) error
	DeleteUserByID(id int) (int, error)
	AuthUser(email string, password string) (*authProto.GeneratedTokens, int, error)
	HashPassword(password string, rounds int) (string, error)
	CheckPasswordHash(password string, hash string) bool
	CheckInputRole(role string) error
	ParseToken(token string) (*authProto.UserRole, error)
	CheckRole(neededRoles []string, givenRole string) error
	CheckRights(neededPerms []string, givenPerms string) error
	RestorePassword(restore *model.RestorePassword) error
}

type Service struct {
	AppUser
}

func NewService(rep *repository.Repository, grpcCli *grpcClient.GRPCClient, logger logging.Logger) *Service {
	return &Service{
		AppUser: NewUserService(*rep, grpcCli, logger),
	}
}
