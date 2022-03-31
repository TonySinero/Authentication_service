package service

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	authProto "stlab.itechart-group.com/go/food_delivery/authentication_service/GRPC"
)

func (u *UserService) AuthUser(email string, password string) (*authProto.GeneratedTokens, int, error) {
	userDb, err := u.repo.AppUser.GetUserByEmail(email)
	if err != nil {
		return nil, 0, err
	}
	if userDb.Deleted {
		u.logger.Errorf("this user (id = %d) is deactivated", userDb.ID)
		return nil, 0, fmt.Errorf("this user (id = %d) is deactivated", userDb.ID)
	}
	if u.CheckPasswordHash(password, userDb.Password) {
		tokens, err := u.grpcCli.TokenGenerationByUserId(context.Background(), &authProto.User{
			UserId: int32(userDb.ID),
			Role:   userDb.Role,
		})
		if err != nil {
			u.logger.Errorf("TokenGenerationByUserId:%s", err)
			return nil, 0, fmt.Errorf("TokenGenerationByUserId:%w", err)
		}
		return tokens, userDb.ID, nil
	} else {
		u.logger.Warn("AuthUser: wrong email or password entered")
		return nil, 0, fmt.Errorf("wrong email or password entered")
	}
}

// HashPassword from string
// bcrypt.DefaultCost = 10
func (u *UserService) HashPassword(password string, rounds int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), rounds)
	return string(bytes), err
}

// CheckPasswordHash compare encrypt
func (u *UserService) CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
