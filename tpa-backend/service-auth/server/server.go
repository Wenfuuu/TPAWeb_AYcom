package server

import (
	"context"

	authpb "github.com/Wenfuuu/tpa-web/pkg/proto/auth"
	userpb "github.com/Wenfuuu/tpa-web/pkg/proto/user"
	"github.com/Wenfuuu/tpa-web/service-auth/handler"
)

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	authHandler *handler.AuthHandler
}

func NewAuthServer(userClient userpb.UserServiceClient) *AuthServer {
	return &AuthServer{
		authHandler: handler.NewAuthHandler(userClient),
	}
}

func (s *AuthServer) Login(ctx context.Context, req *authpb.AuthLoginRequest) (*authpb.AuthLoginResponse, error) {
	return s.authHandler.Login(ctx, req)
}

func (s *AuthServer) Register(ctx context.Context, req *authpb.AuthRegisterRequest) (*authpb.AuthRegisterResponse, error) {
	return s.authHandler.Register(ctx, req)
}

//TODO: authentication & authorization (role)
