package server

import (
	"context"

	threadpb "github.com/Wenfuuu/tpa-web/pkg/proto/thread"
	"github.com/Wenfuuu/tpa-web/service-thread/handler"
)

type ThreadServer struct {
	threadpb.UnimplementedThreadServiceServer
	threadHandler *handler.ThreadHandler
}

func NewThreadServer() *ThreadServer {
	return &ThreadServer{
		threadHandler: handler.NewThreadHandler(),
	}
}

func (s *ThreadServer) GetAllThreads(ctx context.Context, req *threadpb.GetAllThreadsRequest) (*threadpb.GetAllThreadsResponse, error) {
	return s.threadHandler.GetAllThreads(ctx, req)
}

func (s *ThreadServer) PostThread(ctx context.Context, req *threadpb.PostThreadRequest) (*threadpb.ThreadResponse, error) {
	return s.threadHandler.PostThread(ctx, req)
}