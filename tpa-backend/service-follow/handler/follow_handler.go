package handler

import (
	"context"

	followpb "github.com/Wenfuuu/tpa-web/pkg/proto/follow"
	"github.com/Wenfuuu/tpa-web/service-follow/model"
	"github.com/Wenfuuu/tpa-web/service-follow/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FollowHandler struct {
	followRepo *repository.FollowRepository
}

func NewFollowHandler() *FollowHandler {
	return &FollowHandler{
		followRepo: repository.NewFollowRepository(),
	}
}

func (h *FollowHandler) GetFollower(ctx context.Context, req *followpb.GetFollowerRequest) (*followpb.FollowResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "User ID is required")
	}

	modelFollows, err := h.followRepo.GetFollower(req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}

	var pbFollows []*followpb.Follow
	for _, f := range modelFollows {
		pbFollows = append(pbFollows, h.followToProto(&f))
	}

	return &followpb.FollowResponse{
		Follows: pbFollows,
	}, nil
}

func (h *FollowHandler) GetFollowing(ctx context.Context, req *followpb.GetFollowingRequest) (*followpb.FollowResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "User ID is required")
	}

	modelFollows, err := h.followRepo.GetFollowing(req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}

	var pbFollows []*followpb.Follow
	for _, f := range modelFollows {
		pbFollows = append(pbFollows, h.followToProto(&f))
	}

	return &followpb.FollowResponse{
		Follows: pbFollows,
	}, nil
}

func (h *FollowHandler) followToProto(follow *model.Follow) *followpb.Follow {
	return &followpb.Follow{
		FollowerId: follow.FollowerID.String(),
		FollowingId: follow.FollowingID.String(),
	}
}