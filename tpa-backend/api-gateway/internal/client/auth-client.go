package client

import (
	"log"

	authpb "github.com/Wenfuuu/tpa-web/pkg/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAuthClient(authServiceAddr string) (authpb.AuthServiceClient, *grpc.ClientConn) {
	authConn, err := grpc.NewClient(authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
	
	return authpb.NewAuthServiceClient(authConn), authConn
}