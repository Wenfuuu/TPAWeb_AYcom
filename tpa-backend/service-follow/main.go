package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Wenfuuu/tpa-web/service-follow/database"
	"github.com/Wenfuuu/tpa-web/service-follow/server"
	"github.com/Wenfuuu/tpa-web/service-follow/utils"

	followpb "github.com/Wenfuuu/tpa-web/pkg/proto/follow"
	"google.golang.org/grpc"
)

func main() {

	db := database.InitDB()
	database.Migrate(db)
	database.Seed(db)

	port := utils.GetDotEnv("GRPC_PORT")
	if port == "" {
		port = "50055"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	followServer := server.NewFollowServer()
	followpb.RegisterFollowServiceServer(grpcServer, followServer)

	log.Printf("Follow service listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}