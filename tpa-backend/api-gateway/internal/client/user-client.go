package client

import (
	"log"

	// pb "github.com/Wenfuuu/tpa-web/service-user/proto"
	pb "github.com/Wenfuuu/tpa-web/pkg/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserClient(userServiceAddr string) (pb.UserServiceClient) {
	userConn, err := grpc.NewClient(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	
	return pb.NewUserServiceClient(userConn)
}