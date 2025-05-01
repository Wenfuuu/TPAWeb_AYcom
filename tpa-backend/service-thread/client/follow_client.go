package client

import (
	"log"

	followpb "github.com/Wenfuuu/tpa-web/pkg/proto/follow"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewFollowClient(followServiceAddr string) (followpb.FollowServiceClient, *grpc.ClientConn) {
	followConn, err := grpc.NewClient(followServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to follow service: %v", err)
	}

	return followpb.NewFollowServiceClient(followConn), followConn
}