package server

import (
	"context"

	threadpb "github.com/Wenfuuu/tpa-web/pkg/proto/thread"
	userpb "github.com/Wenfuuu/tpa-web/pkg/proto/user"
	"github.com/Wenfuuu/tpa-web/service-thread/handler"
)

type ThreadServer struct {
	threadpb.UnimplementedThreadServiceServer
	threadHandler *handler.ThreadHandler
}

func NewThreadServer(userClient userpb.UserServiceClient) *ThreadServer {
	return &ThreadServer{
		threadHandler: handler.NewThreadHandler(userClient),
	}
}

func (s *ThreadServer) GetTopHashtags(ctx context.Context, req *threadpb.GetTopHashtagsRequest) (*threadpb.HashtagCountResponse, error) {
	return s.threadHandler.GetTopHashtags(ctx, req)
}

func (s *ThreadServer) GetReplies(ctx context.Context, req *threadpb.GetRepliesRequest) (*threadpb.GetAllThreadsResponse, error) {
	return s.threadHandler.GetReplies(ctx, req)
}

func (s *ThreadServer) GetThreadById(ctx context.Context, req *threadpb.GetThreadByIdRequest) (*threadpb.GetThreadResponse, error) {
	return s.threadHandler.GetThreadById(ctx, req)
}

func (s *ThreadServer) GetAllThreads(ctx context.Context, req *threadpb.GetAllThreadsRequest) (*threadpb.GetAllThreadsResponse, error) {
	return s.threadHandler.GetAllThreads(ctx, req)
}

func (s *ThreadServer) GetUserThreads(ctx context.Context, req *threadpb.GetUserThreadsRequest) (*threadpb.GetAllThreadsResponse, error) {
	return s.threadHandler.GetUserThreads(ctx, req)
}

func (s *ThreadServer) GetUserReplies(ctx context.Context, req *threadpb.GetUserRepliesRequest) (*threadpb.GetAllThreadsResponse, error) {
	return s.threadHandler.GetUserReplies(ctx, req)
}

func (s *ThreadServer) GetTopCommunityThreads(ctx context.Context, req *threadpb.GetTopCommunityThreadsRequest) (*threadpb.GetAllThreadsResponse, error) {
	return s.threadHandler.GetTopCommunityThreads(ctx, req)
}

func (s *ThreadServer) GetCommunityThreads(ctx context.Context, req *threadpb.GetCommunityThreadsRequest) (*threadpb.GetAllThreadsResponse, error) {
	return s.threadHandler.GetCommunityThreads(ctx, req)
}

func (s *ThreadServer) GetFollowingThreads(ctx context.Context, req *threadpb.GetFollowingThreadsRequest) (*threadpb.GetAllThreadsResponse, error) {
	return s.threadHandler.GetFollowingThreads(ctx, req)
}

func (s *ThreadServer) GetBookmark(ctx context.Context, req *threadpb.GetBookmarkRequest) (*threadpb.GetAllThreadsResponse, error) {
	return s.threadHandler.GetBookmark(ctx, req)
}

func (s *ThreadServer) GetRepostThreads(ctx context.Context, req *threadpb.GetRepostThreadsRequest) (*threadpb.GetAllThreadsResponse, error) {
	return s.threadHandler.GetRepostThreads(ctx, req)
}

func (s *ThreadServer) GetLikedThreads(ctx context.Context, req *threadpb.GetLikedThreadsRequest) (*threadpb.GetAllThreadsResponse, error) {
	return s.threadHandler.GetLikedThreads(ctx, req)
}

func (s *ThreadServer) GetBookmarkThreads(ctx context.Context, req *threadpb.GetBookmarkThreadsRequest) (*threadpb.GetAllThreadsResponse, error) {
	return s.threadHandler.GetBookmarkThreads(ctx, req)
}

func (s *ThreadServer) PostThread(ctx context.Context, req *threadpb.PostThreadRequest) (*threadpb.ThreadResponse, error) {
	return s.threadHandler.PostThread(ctx, req)
}

func (s *ThreadServer) ReplyThread(ctx context.Context, req *threadpb.ReplyThreadRequest) (*threadpb.ThreadResponse, error) {
	return s.threadHandler.ReplyThread(ctx, req)
}

func (s *ThreadServer) RepostThread(ctx context.Context, req *threadpb.RepostThreadRequest) (*threadpb.ThreadResponse, error) {
	return s.threadHandler.RepostThread(ctx, req)
}

func (s *ThreadServer) UnrepostThread(ctx context.Context, req *threadpb.UnrepostThreadRequest) (*threadpb.ActionResponse, error) {
	return s.threadHandler.UnrepostThread(ctx, req)
}

func (s *ThreadServer) GetCreateCommunity(ctx context.Context, req *threadpb.GetCreateCommunityRequest) (*threadpb.GetCreateCommunityResponse, error) {
	return s.threadHandler.GetCreateCommunity(ctx, req)
}

func (s *ThreadServer) GetJoinCommunity(ctx context.Context, req *threadpb.GetJoinCommunityRequest) (*threadpb.GetJoinCommunityResponse, error) {
	return s.threadHandler.GetJoinCommunity(ctx, req)
}

func (s *ThreadServer) RequestCreateCommunity(ctx context.Context, req *threadpb.RequestCreateCommunityRequest) (*threadpb.ActionResponse, error) {
	return s.threadHandler.RequestCreateCommunity(ctx, req)
}

func (s *ThreadServer) RequestJoinCommunity(ctx context.Context, req *threadpb.RequestJoinCommunityRequest) (*threadpb.ActionResponse, error) {
	return s.threadHandler.RequestJoinCommunity(ctx, req)
}

func (s *ThreadServer) AcceptJoinCommunity(ctx context.Context, req *threadpb.AcceptJoinCommunityRequest) (*threadpb.ActionResponse, error) {
	return s.threadHandler.AcceptJoinCommunity(ctx, req)
}

func (s *ThreadServer) RejectJoinCommunity(ctx context.Context, req *threadpb.RejectJoinCommunityRequest) (*threadpb.ActionResponse, error) {
	return s.threadHandler.RejectJoinCommunity(ctx, req)
}

func (s *ThreadServer) LikeThread(ctx context.Context, req *threadpb.LikeThreadRequest) (*threadpb.ActionResponse, error) {
	return s.threadHandler.LikeThread(ctx, req)
}

func (s *ThreadServer) UnlikeThread(ctx context.Context, req *threadpb.UnlikeThreadRequest) (*threadpb.ActionResponse, error) {
	return s.threadHandler.UnlikeThread(ctx, req)
}

func (s *ThreadServer) BookmarkThread(ctx context.Context, req *threadpb.BookmarkThreadRequest) (*threadpb.ActionResponse, error) {
	return s.threadHandler.BookmarkThread(ctx, req)
}

func (s *ThreadServer) UnbookmarkThread(ctx context.Context, req *threadpb.UnbookmarkThreadRequest) (*threadpb.ActionResponse, error) {
	return s.threadHandler.UnbookmarkThread(ctx, req)
}

func (s *ThreadServer) GetAllThreadCategories(ctx context.Context, req *threadpb.GetAllThreadCategoriesRequest) (*threadpb.GetAllThreadCategoriesResponse, error) {
	return s.threadHandler.GetAllThreadCategories(ctx, req)
}

func (s *ThreadServer) GetCommunityById(ctx context.Context, req *threadpb.GetCommunityByIdRequest) (*threadpb.GetCommunityResponse, error) {
	return s.threadHandler.GetCommunityById(ctx, req)
}

func (s *ThreadServer) GetCommunityMembers(ctx context.Context, req *threadpb.GetCommunityMembersRequest) (*threadpb.CommunityMembersResponse, error) {
	return s.threadHandler.GetCommunityMembers(ctx, req)
}

func (s *ThreadServer) GetTopCommunityMembers(ctx context.Context, req *threadpb.GetTopCommunityMembersRequest) (*threadpb.CommunityMembersResponse, error) {
	return s.threadHandler.GetTopCommunityMembers(ctx, req)
}

func (s *ThreadServer) GetAllCommunity(ctx context.Context, req *threadpb.GetAllCommunitiesRequest) (*threadpb.GetAllCommunitiesResponse, error) {
	return s.threadHandler.GetAllCommunity(ctx, req)
}

func (s *ThreadServer) GetPendingCommunity(ctx context.Context, req *threadpb.GetPendingCommunityRequest) (*threadpb.GetAllCommunitiesResponse, error) {
	return s.threadHandler.GetPendingCommunity(ctx, req)
}

func (s *ThreadServer) GetJoinedCommunity(ctx context.Context, req *threadpb.GetJoinedCommunityRequest) (*threadpb.GetAllCommunitiesResponse, error) {
	return s.threadHandler.GetJoinedCommunity(ctx, req)
}

func (s *ThreadServer) GetCommunityPagination(ctx context.Context, req *threadpb.GetCommunityPaginationRequest) (*threadpb.GetAllCommunitiesResponse, error) {
	return s.threadHandler.GetCommunityPagination(ctx, req)
}

func (s *ThreadServer) GetAllCommunityCategories(ctx context.Context, req *threadpb.GetAllCommunityCategoriesRequest) (*threadpb.GetAllCommunityCategoriesResponse, error) {
	return s.threadHandler.GetAllCommunityCategories(ctx, req)
}

func (s *ThreadServer) GetThreadMedia(ctx context.Context, req *threadpb.GetThreadMediaRequest) (*threadpb.GetAllThreadFileResponse, error) {
	return s.threadHandler.GetThreadMedia(ctx, req)
}

func (s *ThreadServer) GetCommunityMedia(ctx context.Context, req *threadpb.GetCommunityMediaRequest) (*threadpb.GetAllThreadFileResponse, error) {
	return s.threadHandler.GetCommunityMedia(ctx, req)
}

func (s *ThreadServer) GetAllThreadFile(ctx context.Context, req *threadpb.GetAllThreadFileRequest) (*threadpb.GetAllThreadFileResponse, error) {
	return s.threadHandler.GetAllThreadFile(ctx, req)
}