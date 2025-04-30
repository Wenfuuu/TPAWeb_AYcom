package handler

import (
	"context"
	"time"

	threadpb "github.com/Wenfuuu/tpa-web/pkg/proto/thread"
	"github.com/Wenfuuu/tpa-web/service-thread/model"
	"github.com/Wenfuuu/tpa-web/service-thread/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ThreadHandler struct {
	threadRepo *repository.ThreadRepository
}

func NewThreadHandler() *ThreadHandler {
	return &ThreadHandler{
		threadRepo: repository.NewThreadRepository(),
	}
}

func (h *ThreadHandler) GetAllThreads(ctx context.Context, req *threadpb.GetAllThreadsRequest) (*threadpb.GetAllThreadsResponse, error) {
	modelThreads, err := h.threadRepo.GetAllThreads()
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}

	var pbThreads []*threadpb.Thread
	for _, t := range modelThreads {
		pbThreads = append(pbThreads, h.threadToProto(&t))
	}

	return &threadpb.GetAllThreadsResponse{
		Threads: pbThreads,
	}, nil
}

func (h *ThreadHandler) threadToProto(thread *model.Thread) *threadpb.Thread {
	var scheduledAt string
	if thread.ScheduledAt != nil {
		scheduledAt = thread.ScheduledAt.Format(time.RFC3339)
	}


	return &threadpb.Thread{
		Id:        		thread.ID.String(),
		UserId: 		thread.UserID.String(),
		CategoryId: 	thread.CategoryID.String(),
		CommunityId: 	thread.CommunityID.String(),
		Content: 		thread.Content,
		ReplyPermission:thread.ReplyPermission,
		IsAdvertisement:thread.IsAdvertisement,
		Status: 		thread.Status,
		ScheduledAt: 	scheduledAt,
		CreatedAt: 		thread.CreatedAt.Format(time.RFC3339),
	}
}