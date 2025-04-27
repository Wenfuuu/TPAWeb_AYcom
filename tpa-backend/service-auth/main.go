package main

import (
	"fmt"
	"log"
	"net"

	authpb "github.com/Wenfuuu/tpa-web/pkg/proto/auth"
	"github.com/Wenfuuu/tpa-web/service-auth/client"
	"github.com/Wenfuuu/tpa-web/service-auth/server"
	"github.com/Wenfuuu/tpa-web/service-auth/utils"
	"google.golang.org/grpc"
)

func main() {

	// user service
	userServiceAddr := utils.GetDotEnv("USER_SERVICE_ADDR")
	userClient, userConn := client.NewUserClient(userServiceAddr)
	defer userConn.Close()

	port := utils.GetDotEnv("GRPC_PORT")
	if port == "" {
		port = "50052"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	authServer := server.NewAuthServer(userClient)
	authpb.RegisterAuthServiceServer(grpcServer, authServer)

	log.Printf("User service listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}