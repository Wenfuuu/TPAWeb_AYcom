// project/backend/service-user/server/server.go
package server

import (
	"context"

	"github.com/Wenfuuu/tpa-web/service-user/handler"
	pb "github.com/Wenfuuu/tpa-web/service-user/proto"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	userHandler *handler.UserHandler
}

func NewUserServer() *UserServer {
	return &UserServer{
		userHandler: handler.NewUserHandler(),
	}
}

func (s *UserServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.UserResponse, error) {
	return s.userHandler.Register(ctx, req)
}

func (s *UserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return s.userHandler.Login(ctx, req)
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	return s.userHandler.GetUser(ctx, req)
}

func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	return s.userHandler.UpdateUser(ctx, req)
}

func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	return s.userHandler.DeleteUser(ctx, req)
}

func (s *UserServer) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	return s.userHandler.GetAllUsers(ctx, req)
}