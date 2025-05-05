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

func (s *UserServer) GetFollower(ctx context.Context, req *userpb.GetFollowerRequest) (*userpb.FollowResponse, error) {
	return s.userHandler.GetFollower(ctx, req)
}

func (s *UserServer) GetFollowing(ctx context.Context, req *userpb.GetFollowingRequest) (*userpb.FollowResponse, error) {
	return s.userHandler.GetFollowing(ctx, req)
}

func (s *UserServer) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.UserResponse, error) {
	return s.userHandler.Register(ctx, req)
}

func (s *UserServer) RegisterPremium(ctx context.Context, req *userpb.RegisterPremiumRequest) (*userpb.RegisterPremiumResponse, error) {
	return s.userHandler.RegisterPremium(ctx, req)
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

func (s *UserServer) GetWhoToFollow(ctx context.Context, req *userpb.GetWhoToFollowRequest) (*userpb.GetAllUsersResponse, error) {
	return s.userHandler.GetWhoToFollow(ctx, req)
}

func (s *UserServer) UploadImage(ctx context.Context, req *userpb.UploadImageRequest) (*userpb.UploadImageResponse, error) {
	return s.userHandler.UploadImage(ctx, req)
}

func (s *UserServer) UploadVideo(ctx context.Context, req *userpb.UploadVideoRequest) (*userpb.UploadVideoResponse, error) {
	return s.userHandler.UploadVideo(ctx, req)
}