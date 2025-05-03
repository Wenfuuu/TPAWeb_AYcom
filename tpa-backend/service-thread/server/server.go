package server

import (
	"context"

	threadpb "github.com/Wenfuuu/tpa-web/pkg/proto/thread"
	followpb "github.com/Wenfuuu/tpa-web/pkg/proto/follow"
	"github.com/Wenfuuu/tpa-web/service-thread/handler"
)

type ThreadServer struct {
	threadpb.UnimplementedThreadServiceServer
	threadHandler *handler.ThreadHandler
}

func NewThreadServer(followClient followpb.FollowServiceClient) *ThreadServer {
	return &ThreadServer{
		threadHandler: handler.NewThreadHandler(followClient),
	}
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

func (s *ThreadServer) GetFollowingThreads(ctx context.Context, req *threadpb.GetFollowingThreadsRequest) (*threadpb.GetAllThreadsResponse, error) {
	return s.threadHandler.GetFollowingThreads(ctx, req)
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

func (s *ThreadServer) GetAllCommunity(ctx context.Context, req *threadpb.GetAllCommunitiesRequest) (*threadpb.GetAllCommunitiesResponse, error) {
	return s.threadHandler.GetAllCommunity(ctx, req)
}

func (s *ThreadServer) GetAllThreadFile(ctx context.Context, req *threadpb.GetAllThreadFileRequest) (*threadpb.GetAllThreadFileResponse, error) {
	return s.threadHandler.GetAllThreadFile(ctx, req)
}