package service

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	authProto "stlab.itechart-group.com/go/food_delivery/authentication_service/GRPC"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/GRPC/grpcClient"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/mail"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/model"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/pkg/logging"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/repository"
	"strings"
	"time"
)

type UserService struct {
	repo    repository.Repository
	logger  logging.Logger
	grpcCli *grpcClient.GRPCClient
}

func NewUserService(repo repository.Repository, grpcCli *grpcClient.GRPCClient, logger logging.Logger) *UserService {
	return &UserService{repo: repo, grpcCli: grpcCli, logger: logger}
}

func (u *UserService) GetUser(id int) (*model.ResponseUser, error) {
	user, err := u.repo.AppUser.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) GetUsers(page int, limit int, filters *model.RequestFilters) ([]model.ResponseUser, int, error) {
	if filters.Role != "" {
		users, pages, err := u.repo.AppUser.GetUserByRoleFilter(page, limit, filters)
		if err != nil {
			return nil, 0, err
		}
		return users, pages, nil
	} else if filters.FilterData {
		if filters.EndTime.Unix() < filters.StartTime.Unix() {
			filters.EndTime.Time = filters.StartTime.Time
		}
		users, pages, err := u.repo.AppUser.GetUserByDataFilter(page, limit, filters)
		if err != nil {
			return nil, 0, err
		}
		return users, pages, nil
	} else {
		users, pages, err := u.repo.AppUser.GetUserAll(page, limit)
		if err != nil {
			return nil, 0, err
		}
		return users, pages, nil
	}
}

func (u *UserService) CreateCustomer(user *model.CreateCustomer) (*authProto.GeneratedTokens, int, error) {
	if user.Password == "" {
		user.Password = GeneratePassword()
	}
	pas := user.Password
	hash, err := u.HashPassword(user.Password, bcrypt.DefaultCost)
	if err != nil {
		u.logger.Errorf("CreateUser: can not generate hash from password:%s", err)
		return nil, 0, fmt.Errorf("createUser: can not generate hash from password:%w", err)
	}
	user.Password = hash
	id, err := u.repo.AppUser.CreateCustomer(user)
	if err != nil {
		return nil, 0, err
	}
	go mail.SendEmail(u.logger, &model.Post{
		Email:    user.Email,
		Password: pas,
	})
	_, err = u.grpcCli.BindUserAndRole(context.Background(), &authProto.User{
		UserId: int32(id),
		Role:   "Authorized Customer",
	})
	if err != nil {
		u.logger.Errorf("BindUserAndRole:%s", err)
		return nil, id, fmt.Errorf("bindUserAndRole:%w", err)
	}
	tokens, err := u.grpcCli.TokenGenerationByUserId(context.Background(), &authProto.User{
		UserId: int32(id),
		Role:   "Authorized Customer",
	})
	if err != nil {
		u.logger.Errorf("tokenGenerationByUserId:%s", err)
		return nil, 0, fmt.Errorf("tokenGenerationByUserId:%w", err)
	}
	return tokens, id, nil
}

func (u *UserService) CreateStaff(user *model.CreateStaff) (int, error) {
	if user.Password == "" {
		user.Password = GeneratePassword()
	}
	pas := user.Password
	hash, err := u.HashPassword(user.Password, bcrypt.DefaultCost)
	if err != nil {
		u.logger.Errorf("CreateStaff: can not generate hash from password:%s", err)
		return 0, fmt.Errorf("CreateStaff: can not generate hash from password:%w", err)
	}
	user.Password = hash
	id, err := u.repo.AppUser.CreateStaff(user)
	if err != nil {
		return 0, err
	}
	go mail.SendEmail(u.logger, &model.Post{
		Email:    user.Email,
		Password: pas,
	})
	_, err = u.grpcCli.BindUserAndRole(context.Background(), &authProto.User{
		UserId: int32(id),
		Role:   user.Role,
	})
	if err != nil {
		u.logger.Errorf("BindUserAndRole:%s", err)
		return id, fmt.Errorf("bindUserAndRole:%w", err)
	}
	return id, nil
}

func (u *UserService) UpdateUser(user *model.UpdateUser) error {
	userDb, err := u.repo.AppUser.GetUserByEmail(user.Email)
	if err != nil {
		return err
	}
	if u.CheckPasswordHash(user.OldPassword, userDb.Password) {
		newHash, err := u.HashPassword(user.NewPassword, bcrypt.DefaultCost)
		if err != nil {
			u.logger.Errorf("UpdateUser: can not generate hash from password:%s", err)
			return fmt.Errorf("updateUser: can not generate hash from password:%w", err)
		}
		user.NewPassword = newHash
		err = u.repo.AppUser.UpdateUser(user)
		if err != nil {
			return err
		}
		return nil
	} else {
		u.logger.Warn("wrong email or password entered")
		return fmt.Errorf("wrong email or password entered")
	}
}

func (u *UserService) DeleteUserByID(id int) (int, error) {
	userId, err := u.repo.AppUser.DeleteUserByID(id)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (u *UserService) CheckInputRole(role string) error {
	roles, err := u.grpcCli.GetAllRoles(context.Background(), &empty.Empty{})
	if err != nil {
		u.logger.Errorf("CheckInputRole:%s", err)
		return err
	}
	roleSlice := strings.Split(roles.Roles, ",")
	var resultCheck = false
	for _, saveRole := range roleSlice {
		if saveRole == role {
			resultCheck = true
		}
	}
	if resultCheck == false {
		return fmt.Errorf("incorrect role in request")
	} else {
		return nil
	}
}

func (u *UserService) ParseToken(token string) (*authProto.UserRole, error) {
	return u.grpcCli.GetUserWithRights(context.Background(), &authProto.AccessToken{AccessToken: token})
}

func (u *UserService) CheckRole(neededRoles []string, givenRole string) error {
	neededRolesString := strings.Join(neededRoles, ",")
	if !strings.Contains(neededRolesString, givenRole) {
		return fmt.Errorf("not enough rights")
	}
	return nil
}

func (u *UserService) CheckRights(neededPerms []string, givenPerms string) error {
	if neededPerms != nil {
		ok := true
		for _, perm := range neededPerms {
			if !strings.Contains(givenPerms, perm) {
				ok = false
				return fmt.Errorf("not enough rights")
			} else {
				continue
			}
		}
		if ok == true {
			return nil
		}
	}
	return nil
}

func (u *UserService) RestorePassword(restore *model.RestorePassword) error {
	err := u.repo.CheckEmail(restore.Email)
	if err != nil {
		return err
	}
	password := GeneratePassword()
	hash, err := u.HashPassword(password, bcrypt.DefaultCost)
	if err != nil {
		u.logger.Errorf("RestorePassword: can not generate hash from password:%s", err)
		return fmt.Errorf("RestorePassword: can not generate hash from password:%w", err)
	}
	restore.Password = hash
	err = u.repo.AppUser.RestorePassword(restore)
	if err != nil {
		return err
	}
	go mail.SendEmail(u.logger, &model.Post{
		Email:    restore.Email,
		Password: password,
	})
	return nil
}

func GeneratePassword() string {
	rand.Seed(time.Now().UnixNano())
	length := 8 + rand.Intn(7)
	var b strings.Builder
	b.WriteRune(model.PasswordUpper[rand.Intn(len(model.PasswordUpper))])
	b.WriteRune(model.PasswordNumber[rand.Intn(len(model.PasswordNumber))])
	b.WriteRune(model.PasswordLower[rand.Intn(len(model.PasswordLower))])
	b.WriteRune(model.PasswordSpecial[rand.Intn(len(model.PasswordSpecial))])
	for i := 0; i < length-4; i++ {
		b.WriteRune(model.PasswordComposition[rand.Intn(len(model.PasswordComposition))])
	}
	return b.String()
}
