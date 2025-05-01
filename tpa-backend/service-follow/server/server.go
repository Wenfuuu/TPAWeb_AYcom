package server

import (
	"context"

	followpb "github.com/Wenfuuu/tpa-web/pkg/proto/follow"
	"github.com/Wenfuuu/tpa-web/service-follow/handler"
)

type FollowServer struct {
	followpb.UnimplementedFollowServiceServer
	followHandler *handler.FollowHandler
}

func NewFollowServer() *FollowServer {
	return &FollowServer{
		followHandler: handler.NewFollowHandler(),
	}
}

func (s *FollowServer) GetFollower(ctx context.Context, req *followpb.GetFollowerRequest) (*followpb.FollowResponse, error) {
	return s.followHandler.GetFollower(ctx, req)
}

func (s *FollowServer) GetFollowing(ctx context.Context, req *followpb.GetFollowingRequest) (*followpb.FollowResponse, error) {
	return s.followHandler.GetFollowing(ctx, req)
}