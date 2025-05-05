package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	threadpb "github.com/Wenfuuu/tpa-web/pkg/proto/thread"
	userpb "github.com/Wenfuuu/tpa-web/pkg/proto/user"
	"github.com/Wenfuuu/tpa-web/service-thread/model"
	"github.com/Wenfuuu/tpa-web/service-thread/repository"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ThreadHandler struct {
	threadRepo *repository.ThreadRepository
	userClient userpb.UserServiceClient
}

func NewThreadHandler(userClient userpb.UserServiceClient) *ThreadHandler {
	return &ThreadHandler{
		threadRepo: repository.NewThreadRepository(),
		userClient: userClient,
	}
}

func (h *ThreadHandler) GetTopHashtags(ctx context.Context, req *threadpb.GetTopHashtagsRequest) (*threadpb.HashtagCountResponse, error) {
	modelHashtags, err := h.threadRepo.GetTopHashtags();
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}

	var pbHashtags []*threadpb.HashtagCount
	for _, t := range modelHashtags {
		pbHashtags = append(pbHashtags, h.hashtagCountToProto(&t))
	}

	return &threadpb.HashtagCountResponse{
		HashtagCount: pbHashtags,
	}, nil
}

func (h *ThreadHandler) hashtagCountToProto(hashtag *model.HashtagCount) *threadpb.HashtagCount {
	return &threadpb.HashtagCount{
		Id: 	hashtag.ID.String(),
		Name: 	hashtag.Name,
		Count: 	int32(hashtag.Count),
	}
}

func (h *ThreadHandler) GetReplies(ctx context.Context, req *threadpb.GetRepliesRequest) (*threadpb.GetAllThreadsResponse, error) {
	page := int(req.Page)
	modelThreads, err := h.threadRepo.GetReplies(page, req.ThreadId)
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

func (h *ThreadHandler) GetThreadById(ctx context.Context, req *threadpb.GetThreadByIdRequest) (*threadpb.GetThreadResponse, error) {
	modelThread, err := h.threadRepo.GetThreadById(req.ThreadId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}

	return &threadpb.GetThreadResponse{
		Thread: h.threadToProto(modelThread),
	}, nil
}

func (h *ThreadHandler) GetAllThreads(ctx context.Context, req *threadpb.GetAllThreadsRequest) (*threadpb.GetAllThreadsResponse, error) {
	page := int(req.Page)
	modelThreads, err := h.threadRepo.GetAllThreads(page)
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

func (h *ThreadHandler) GetCommunityThreads(ctx context.Context, req *threadpb.GetCommunityThreadsRequest) (*threadpb.GetAllThreadsResponse, error) {
	page := int(req.Page)
	communityId := req.CommunityId

	modelThreads, err := h.threadRepo.GetCommunityThreads(page, communityId)
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

func (h *ThreadHandler) GetFollowingThreads(ctx context.Context, req *threadpb.GetFollowingThreadsRequest) (*threadpb.GetAllThreadsResponse, error) {
	page := int(req.Page)
	userId := req.UserId // the current account id

	// call user service to get following ids
	userResp, err := h.userClient.GetFollowing(ctx, &userpb.GetFollowingRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get followings")
	}

	var followingIds []string
	for _, element := range userResp.Follows {
		if element != nil {
			followingIds = append(followingIds, element.FollowingId)
		}
	}

	followingIds = append(followingIds, userId)
	modelThreads, err := h.threadRepo.GetFollowingThreads(page, followingIds)
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

func (h *ThreadHandler) GetBookmark(ctx context.Context, req *threadpb.GetBookmarkRequest) (*threadpb.GetAllThreadsResponse, error) {
	page := int(req.Page)
	modelThreads, err := h.threadRepo.GetBookmark(page, req.UserId, req.Search)
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

func (h *ThreadHandler) GetRepostThreads(ctx context.Context, req *threadpb.GetRepostThreadsRequest) (*threadpb.GetAllThreadsResponse, error) {
	modelThreads, err := h.threadRepo.GetRepostThreads(req.UserId)
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

func (h *ThreadHandler) GetLikedThreads(ctx context.Context, req *threadpb.GetLikedThreadsRequest) (*threadpb.GetAllThreadsResponse, error) {
	modelThreads, err := h.threadRepo.GetLikedThreads(req.UserId)
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

func (h *ThreadHandler) GetBookmarkThreads(ctx context.Context, req *threadpb.GetBookmarkThreadsRequest) (*threadpb.GetAllThreadsResponse, error) {
	modelThreads, err := h.threadRepo.GetBookmarkThreads(req.UserId)
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

func (h *ThreadHandler) UnrepostThread(ctx context.Context, req *threadpb.UnrepostThreadRequest) (*threadpb.ActionResponse, error) {
	userID, err := uuid.Parse(req.UserId)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %v", err)
    }

    threadID, err := uuid.Parse(req.ThreadId)
    if err != nil {
        return nil, fmt.Errorf("invalid category ID: %v", err)
    }

	repost := &model.Thread{
		UserID:      userID,
		RepostOfID:  &threadID,
	}

	if err := h.threadRepo.UnrepostThread(repost); err != nil {
		return nil, fmt.Errorf("failed to unrepost thread: %v", err)
	}

	return &threadpb.ActionResponse{
		Message: "Success unrepost a thread!",
	}, nil
}

func (h *ThreadHandler) RepostThread(ctx context.Context, req *threadpb.RepostThreadRequest) (*threadpb.ThreadResponse, error) {
	userID, err := uuid.Parse(req.UserId)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %v", err)
    }

	modelThread, err := h.threadRepo.GetThreadById(req.ThreadId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}
	
	thread := model.Thread{
		UserID: userID,
		RepostOfID: &modelThread.ID,
		Content: modelThread.Content,
		ReplyPermission: modelThread.ReplyPermission,
		IsAdvertisement: modelThread.IsAdvertisement,
		ScheduledAt: modelThread.ScheduledAt,
	}

	if err := h.threadRepo.Create(&thread); err != nil {
		return nil, status.Error(codes.Internal, "Failed to repost thread")
	}

	fileUrls, err := h.threadRepo.GetFileById(modelThread.ID.String())// img url from the reposted thread
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}

	var threadFiles []model.ThreadFiles
	for _, url := range fileUrls {
		threadFiles = append(threadFiles, model.ThreadFiles{
			ThreadID: thread.ID,
			File:     url.File,
		})
	}

	if len(threadFiles) > 0 {
		if err = h.threadRepo.InsertFile(threadFiles); err != nil {
			return nil, status.Error(codes.Internal, "Failed to save thread files")
		}
	}

	return h.threadToResponse(&thread), nil
}

func (h *ThreadHandler) ReplyThread(ctx context.Context, req *threadpb.ReplyThreadRequest) (*threadpb.ThreadResponse, error) {
	userID, err := uuid.Parse(req.UserId)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %v", err)
    }

	var replyToId *uuid.UUID
	if req.ReplyToId != "" {
		parsedReply, err := uuid.Parse(req.ReplyToId)
		if err != nil {
			return nil, fmt.Errorf("invalid reply id: %v", err)
		}
		replyToId = &parsedReply
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
		ReplyToID: replyToId,
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

func (h *ThreadHandler) JoinCommunity(ctx context.Context, req *threadpb.JoinCommunityRequest) (*threadpb.ActionResponse, error) {
	userID, err := uuid.Parse(req.UserId)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %v", err)
    }

    communityID, err := uuid.Parse(req.CommunityId)
    if err != nil {
        return nil, fmt.Errorf("invalid community ID: %v", err)
    }

	member := model.CommunityMembers{
		UserID: userID,
		CommunityID: communityID,
	}

	if err := h.threadRepo.JoinCommunity(&member, req.UserId, req.CommunityId); err != nil {
		return nil, status.Error(codes.Internal, "Failed to join community!")
	}

	return &threadpb.ActionResponse{
		Message: "Success join a community!",
	}, nil
}

func (h *ThreadHandler) RequestJoinCommunity(ctx context.Context, req *threadpb.RequestJoinCommunityRequest) (*threadpb.ActionResponse, error) {
	userID, err := uuid.Parse(req.UserId)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %v", err)
    }

    communityID, err := uuid.Parse(req.CommunityId)
    if err != nil {
        return nil, fmt.Errorf("invalid community ID: %v", err)
    }

	request := model.JoinCommunityRequest{
		UserID: userID,
		CommunityID: communityID,
	}

	if err := h.threadRepo.RequestJoinCommunity(&request); err != nil {
		return nil, status.Error(codes.Internal, "Failed to request join community!")
	}

	return &threadpb.ActionResponse{
		Message: "Success request join a community!",
	}, nil
}

func (h *ThreadHandler) LikeThread(ctx context.Context, req *threadpb.LikeThreadRequest) (*threadpb.ActionResponse, error) {
	userID, err := uuid.Parse(req.UserId)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %v", err)
    }

    threadID, err := uuid.Parse(req.ThreadId)
    if err != nil {
        return nil, fmt.Errorf("invalid thread ID: %v", err)
    }

	like := model.Like{
		UserID: userID,
		ThreadID: threadID,
	}

	if err := h.threadRepo.LikeThread(&like); err != nil {
		return nil, status.Error(codes.Internal, "Failed to like thread")
	}

	return &threadpb.ActionResponse{
		Message: "Success like a thread!",
	}, nil
}

func (h *ThreadHandler) UnlikeThread(ctx context.Context, req *threadpb.UnlikeThreadRequest) (*threadpb.ActionResponse, error) {
	userID, err := uuid.Parse(req.UserId)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %v", err)
    }

    threadID, err := uuid.Parse(req.ThreadId)
    if err != nil {
        return nil, fmt.Errorf("invalid category ID: %v", err)
    }

	unlike := model.Like{
		UserID: userID,
		ThreadID: threadID,
	}

	if err := h.threadRepo.UnlikeThread(&unlike); err != nil {
		return nil, status.Error(codes.Internal, "Failed to unlike thread")
	}

	return &threadpb.ActionResponse{
		Message: "Success unlike a thread!",
	}, nil
}

func (h *ThreadHandler) BookmarkThread(ctx context.Context, req *threadpb.BookmarkThreadRequest) (*threadpb.ActionResponse, error) {
	userID, err := uuid.Parse(req.UserId)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %v", err)
    }

    threadID, err := uuid.Parse(req.ThreadId)
    if err != nil {
        return nil, fmt.Errorf("invalid category ID: %v", err)
    }

	bookmark := model.Bookmark{
		UserID: userID,
		ThreadID: threadID,
	}

	if err := h.threadRepo.BookmarkThread(&bookmark); err != nil {
		return nil, status.Error(codes.Internal, "Failed to bookmark thread")
	}

	return &threadpb.ActionResponse{
		Message: "Success bookmark a thread!",
	}, nil
}

func (h *ThreadHandler) UnbookmarkThread(ctx context.Context, req *threadpb.UnbookmarkThreadRequest) (*threadpb.ActionResponse, error) {
	userID, err := uuid.Parse(req.UserId)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %v", err)
    }

    threadID, err := uuid.Parse(req.ThreadId)
    if err != nil {
        return nil, fmt.Errorf("invalid category ID: %v", err)
    }

	unboookmark := model.Bookmark{
		UserID: userID,
		ThreadID: threadID,
	}

	if err := h.threadRepo.UnbookmarkThread(&unboookmark); err != nil {
		return nil, status.Error(codes.Internal, "Failed to unbookmark thread")
	}

	return &threadpb.ActionResponse{
		Message: "Success unbookmark a thread!",
	}, nil
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

func (h *ThreadHandler) GetCommunityById(ctx context.Context, req *threadpb.GetCommunityByIdRequest) (*threadpb.GetCommunityResponse, error) {
	modelCommunity, err := h.threadRepo.GetCommunityById(req.CommunityId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}

	return &threadpb.GetCommunityResponse{
		Community: h.communityToProto(modelCommunity),
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

func (h *ThreadHandler) GetPendingCommunity(ctx context.Context, req *threadpb.GetPendingCommunityRequest) (*threadpb.GetAllCommunitiesResponse, error) {
	page := int(req.Page)
	limit := int(req.Limit)

	modelCommunities, err := h.threadRepo.GetPendingCommunity(req.UserId, page, limit, req.Search, req.CategoryId)
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

func (h *ThreadHandler) GetJoinedCommunity(ctx context.Context, req *threadpb.GetJoinedCommunityRequest) (*threadpb.GetAllCommunitiesResponse, error) {
	page := int(req.Page)
	limit := int(req.Limit)

	modelCommunities, err := h.threadRepo.GetJoinedCommunity(req.UserId, page, limit, req.Search, req.CategoryId)
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

func (h *ThreadHandler) GetCommunityPagination(ctx context.Context, req *threadpb.GetCommunityPaginationRequest) (*threadpb.GetAllCommunitiesResponse, error) {
	page := int(req.Page)
	limit := int(req.Limit)

	modelCommunities, err := h.threadRepo.GetCommunityPagination(page, limit, req.Search, req.CategoryId)
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

func (h *ThreadHandler) GetAllCommunityCategories(ctx context.Context, req *threadpb.GetAllCommunityCategoriesRequest) (*threadpb.GetAllCommunityCategoriesResponse, error) {
	modelCategories, err := h.threadRepo.GetAllCommunityCategories()
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}

	var pbThreads []*threadpb.CommunityCategory
	for _, t := range modelCategories {
		pbThreads = append(pbThreads, h.communityCategoryToProto(&t))
	}

	return &threadpb.GetAllCommunityCategoriesResponse{
		CommunityCategories: pbThreads,
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
		Logo: 			community.Logo,	
		CategoryId: 	community.CategoryID.String(),
		Banner: 		community.Banner,
		Rule: 			community.Rule,
		CreatorId: 		community.CreatorID.String(),
		CreatedAt: 		community.CreatedAt.Format(time.RFC3339),
	}
}

func (h *ThreadHandler) categoryToProto(category *model.ThreadCategory) *threadpb.ThreadCategory {
	return &threadpb.ThreadCategory{
		Id: 	category.ID.String(),
		Name: 	category.Name,
	}
}

func (h *ThreadHandler) communityCategoryToProto(category *model.CommunityCategory) *threadpb.CommunityCategory {
	return &threadpb.CommunityCategory{
		Id: 	category.ID.String(),
		Name: 	category.Name,
	}
}

func (h *ThreadHandler) threadToProto(thread *model.Thread) *threadpb.Thread {
	var scheduledAt string
	if thread.ScheduledAt != nil {
		scheduledAt = thread.ScheduledAt.Format(time.RFC3339)
	}

	var replyToId string
	if thread.ReplyToID != nil {
		replyToId = thread.ReplyToID.String()
	}

	var repostOfId string
	if thread.RepostOfID != nil {
		repostOfId = thread.RepostOfID.String()
	}

	return &threadpb.Thread{
		Id:        		thread.ID.String(),
		UserId: 		thread.UserID.String(),
		CategoryId: 	thread.CategoryID.String(),
		CommunityId: 	thread.CommunityID.String(),
		ReplyToId: 		replyToId,
		RepostOfId: 	repostOfId,
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