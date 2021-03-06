// Code generated by protoc-gen-defaults. DO NOT EDIT.

package mockAuthProto

import (
	"context"
	"github.com/bxcodec/faker"
	"github.com/golang/protobuf/ptypes/empty"
	authProto "stlab.itechart-group.com/go/food_delivery/authentication_service/GRPC"
)

// MockAuthServer is the mock implementation of the AuthServer. Use this to create mock services that
// return random data. Useful in UI Testing.
type MockAuthServer struct {
	authProto.UnimplementedAuthServer
}

// GetUserWithRights is mock implementation of the method GetUserWithRights
func (*MockAuthServer) GetUserWithRights(context.Context, *authProto.AccessToken) (*authProto.UserRole, error) {
	var res authProto.UserRole
	if err := faker.FakeData(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

// BindUserAndRole is mock implementation of the method BindUserAndRole
func (*MockAuthServer) BindUserAndRole(context.Context, *authProto.User) (*authProto.ResultBinding, error) {
	var res authProto.ResultBinding
	if err := faker.FakeData(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

// TokenGenerationByRefresh is mock implementation of the method TokenGenerationByRefresh
func (*MockAuthServer) TokenGenerationByRefresh(context.Context, *authProto.RefreshToken) (*authProto.GeneratedTokens, error) {
	var res authProto.GeneratedTokens
	if err := faker.FakeData(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

// TokenGenerationByUserId is mock implementation of the method TokenGenerationByUserId
func (*MockAuthServer) TokenGenerationByUserId(context.Context, *authProto.User) (*authProto.GeneratedTokens, error) {
	var res authProto.GeneratedTokens
	if err := faker.FakeData(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetAllRoles is mock implementation of the method GetAllRoles
func (*MockAuthServer) GetAllRoles(context.Context, *empty.Empty) (*authProto.Roles, error) {
	var res authProto.Roles
	if err := faker.FakeData(&res); err != nil {
		return nil, err
	}
	return &res, nil
}
