package client

import (
	"log"

	chatpb "github.com/Wenfuuu/tpa-web/pkg/proto/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewChatClient(chatServiceAddr string) (chatpb.ChatServiceClient, *grpc.ClientConn) {
	chatConn, err := grpc.NewClient(chatServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to chat service: %v", err)
	}

	return chatpb.NewChatServiceClient(chatConn), chatConn
}