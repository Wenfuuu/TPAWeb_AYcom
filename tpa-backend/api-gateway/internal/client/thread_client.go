package client

import (
	"log"

	threadpb "github.com/Wenfuuu/tpa-web/pkg/proto/thread"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewThreadClient(threadServiceAddr string) (threadpb.ThreadServiceClient, *grpc.ClientConn) {
	threadConn, err := grpc.NewClient(threadServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to thread service: %v", err)
	}
	
	return threadpb.NewThreadServiceClient(threadConn), threadConn
}