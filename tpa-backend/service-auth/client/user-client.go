package client

import (
	"log"

	userpb "github.com/Wenfuuu/tpa-web/pkg/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserClient(userServiceAddr string) (userpb.UserServiceClient, *grpc.ClientConn) {
	userConn, err := grpc.NewClient(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	
	return userpb.NewUserServiceClient(userConn), userConn
}