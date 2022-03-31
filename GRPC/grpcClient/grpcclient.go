package grpcClient

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	authProto "stlab.itechart-group.com/go/food_delivery/authentication_service/GRPC"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/pkg/logging"
)

var logger = logging.GetLogger()

type GRPCClient struct {
	cli authProto.AuthClient
}

func NewGRPCClient(host string) *GRPCClient {
	Target := fmt.Sprintf("%s:8090", host)
	conn, err := grpc.Dial(Target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("NewGRPCClient, Dial:%s", err)
	}
	cli := authProto.NewAuthClient(conn)
	return &GRPCClient{cli: cli}
}

func (c *GRPCClient) GetUserWithRights(ctx context.Context, in *authProto.AccessToken, opts ...grpc.CallOption) (*authProto.UserRole, error) {
	return c.cli.GetUserWithRights(ctx, in)
}

func (c *GRPCClient) BindUserAndRole(ctx context.Context, in *authProto.User, opts ...grpc.CallOption) (*authProto.ResultBinding, error) {
	return c.cli.BindUserAndRole(ctx, in)
}

func (c *GRPCClient) TokenGenerationByRefresh(ctx context.Context, in *authProto.RefreshToken, opts ...grpc.CallOption) (*authProto.GeneratedTokens, error) {
	return nil, nil
}

func (c *GRPCClient) TokenGenerationByUserId(ctx context.Context, in *authProto.User, opts ...grpc.CallOption) (*authProto.GeneratedTokens, error) {
	return c.cli.TokenGenerationByUserId(ctx, in)
}

func (c *GRPCClient) GetAllRoles(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*authProto.Roles, error) {
	return c.cli.GetAllRoles(ctx, in)
}
