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

func (s *ThreadServer) GetAllThreads(ctx context.Context, req *threadpb.GetAllThreadsRequest) (*threadpb.GetAllThreadsResponse, error) {
	return s.threadHandler.GetAllThreads(ctx, req)
}

func (s *ThreadServer) GetFollowingThreads(ctx context.Context, req *threadpb.GetFollowingThreadsRequest) (*threadpb.GetAllThreadsResponse, error) {
	return s.threadHandler.GetFollowingThreads(ctx, req)
}

func (s *ThreadServer) PostThread(ctx context.Context, req *threadpb.PostThreadRequest) (*threadpb.ThreadResponse, error) {
	return s.threadHandler.PostThread(ctx, req)
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