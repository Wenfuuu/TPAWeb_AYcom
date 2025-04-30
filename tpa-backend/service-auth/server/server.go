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
func (s *AuthServer) CheckAuth(ctx context.Context, req *authpb.CheckAuthRequest) (*authpb.CheckAuthResponse, error) {
	return s.authHandler.CheckAuth(ctx, req)
}

func (s *AuthServer) AuthorizeAdmin(ctx context.Context, req *authpb.AuthorizeAdminRequest) (*authpb.AuthorizeAdminResponse, error) {
	return s.authHandler.AuthorizeAdmin(ctx, req)
}

func (s *AuthServer) GetAuthUser(ctx context.Context, req *authpb.GetAuthUserRequest) (*authpb.AuthLoginResponse, error) {
	return s.authHandler.GetAuthUser(ctx, req)
}

func (s *AuthServer) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.AuthLoginResponse, error) {
	return s.authHandler.RefreshToken(ctx, req)
}

func (s *AuthServer) SendVerificationCode(ctx context.Context, req *authpb.SendVerificationCodeRequest) (*authpb.SendVerificationCodeResponse, error) {
	return s.authHandler.SendVerificationCode(ctx, req)
}

func (s *AuthServer) VerifyCode(ctx context.Context, req *authpb.VerifyCodeRequest) (*authpb.VerifyCodeResponse, error) {
	return s.authHandler.VerifyCode(ctx, req)
}