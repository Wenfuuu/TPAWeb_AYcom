// project/backend/service-user/server/server.go
package server

import (
	"context"

	"github.com/Wenfuuu/tpa-web/service-user/handler"
	// pb "github.com/Wenfuuu/tpa-web/service-user/proto"
	userpb "github.com/Wenfuuu/tpa-web/pkg/proto/user"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
	userHandler *handler.UserHandler
}

func NewUserServer() *UserServer {
	return &UserServer{
		userHandler: handler.NewUserHandler(),
	}
}

func (s *UserServer) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.UserResponse, error) {
	return s.userHandler.Register(ctx, req)
}

func(s *UserServer) GetUserByEmail(ctx context.Context, req *userpb.GetUserByEmailRequest) (*userpb.UserResponse, error) {
	return s.userHandler.GetUserByEmail(ctx, req)
}

func (s *UserServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	return s.userHandler.GetUser(ctx, req)
}

func (s *UserServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {
	return s.userHandler.UpdateUser(ctx, req)
}

func (s *UserServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	return s.userHandler.DeleteUser(ctx, req)
}

func (s *UserServer) GetAllUsers(ctx context.Context, req *userpb.GetAllUsersRequest) (*userpb.GetAllUsersResponse, error) {
	return s.userHandler.GetAllUsers(ctx, req)
}