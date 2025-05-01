package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	threadpb "github.com/Wenfuuu/tpa-web/pkg/proto/thread"
	"github.com/Wenfuuu/tpa-web/service-thread/model"
	"github.com/Wenfuuu/tpa-web/service-thread/repository"
	"github.com/google/uuid"
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

func (h *ThreadHandler) PostThread(ctx context.Context, req *threadpb.PostThreadRequest) (*threadpb.ThreadResponse, error) {
	log.Println(req.UserId)
	log.Println(req.CategoryId)
	log.Println(req.CommunityId)

	userID, err := uuid.Parse(req.UserId)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %v", err)
    }

    categoryID, err := uuid.Parse(req.CategoryId)
    if err != nil {
        return nil, fmt.Errorf("invalid category ID: %v", err)
    }

	var communityID uuid.UUID
	if(req.CommunityId != ""){
		communityID, err = uuid.Parse(req.CommunityId)
		if err != nil {
			return nil, fmt.Errorf("invalid community ID: %v", err)
		}
	}

	var scheduledAt *time.Time
	if req.ScheduledAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, req.ScheduledAt)
		if err != nil {
			return nil, fmt.Errorf("invalid scheduled_at: %v", err)
		}
		scheduledAt = &parsedTime
	}

	thread := model.Thread{
		UserID: userID,
		CategoryID: categoryID,
		CommunityID: communityID,
		Content: req.Content,
		ReplyPermission: req.ReplyPermission,
		IsAdvertisement: req.IsAdvertisement,
		ScheduledAt: scheduledAt,
	}

	if err := h.threadRepo.Create(&thread); err != nil {
		return nil, status.Error(codes.Internal, "Failed to create thread")
	}

	var threadFiles []model.ThreadFiles
	for _, url := range req.ImageUrls {
		threadFiles = append(threadFiles, model.ThreadFiles{
			ThreadID: thread.ID,
			File:     url,
		})
	}
	
	if len(threadFiles) > 0 {
		if err = h.threadRepo.InsertFile(threadFiles); err != nil {
			return nil, status.Error(codes.Internal, "Failed to save thread files")
		}
	}

	return h.threadToResponse(&thread), nil
}

func (h *ThreadHandler) GetAllThreadCategories(ctx context.Context, req *threadpb.GetAllThreadCategoriesRequest) (*threadpb.GetAllThreadCategoriesResponse, error) {
	modelCategories, err := h.threadRepo.GetAllThreadCategories()
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}

	var pbThreads []*threadpb.ThreadCategory
	for _, t := range modelCategories {
		pbThreads = append(pbThreads, h.categoryToProto(&t))
	}

	return &threadpb.GetAllThreadCategoriesResponse{
		ThreadCategories: pbThreads,
	}, nil
}

func (h *ThreadHandler) GetAllCommunity(ctx context.Context, req *threadpb.GetAllCommunitiesRequest) (*threadpb.GetAllCommunitiesResponse, error) {
	modelCommunities, err := h.threadRepo.GetAllCommunity()
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}

	var pbThreads []*threadpb.Community
	for _, t := range modelCommunities {
		pbThreads = append(pbThreads, h.communityToProto(&t))
	}

	return &threadpb.GetAllCommunitiesResponse{
		Communities: pbThreads,
	}, nil
}

func (h *ThreadHandler) GetAllThreadFile(ctx context.Context, req *threadpb.GetAllThreadFileRequest) (*threadpb.GetAllThreadFileResponse, error) {
	modelFiles, err := h.threadRepo.GetAllThreadFile()
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}

	var pbThreads []*threadpb.ThreadFile
	for _, t := range modelFiles {
		pbThreads = append(pbThreads, h.fileToProto(&t))
	}

	return &threadpb.GetAllThreadFileResponse{
		ThreadFiles: pbThreads,
	}, nil
}

func (h *ThreadHandler) fileToProto(file *model.ThreadFiles) *threadpb.ThreadFile {
	return &threadpb.ThreadFile{
		Id: 		file.ID.String(),
		ThreadId: 	file.ThreadID.String(),
		File: 		file.File,
	}
}

func (h *ThreadHandler) communityToProto(community *model.Community) *threadpb.Community {
	return &threadpb.Community{
		Id: 			community.ID.String(),
		Name: 			community.Name,
		Description: 	community.Description,
	}
}

func (h *ThreadHandler) categoryToProto(category *model.ThreadCategory) *threadpb.ThreadCategory {
	return &threadpb.ThreadCategory{
		Id: 	category.ID.String(),
		Name: 	category.Name,
	}
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

func (h *ThreadHandler) threadToResponse(thread *model.Thread) *threadpb.ThreadResponse {
	return &threadpb.ThreadResponse{
		Thread: h.threadToProto(thread),
	}
}