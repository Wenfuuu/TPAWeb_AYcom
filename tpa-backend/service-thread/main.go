package main

import (
	"github.com/Wenfuuu/tpa-web/service-thread/client"
	"github.com/Wenfuuu/tpa-web/service-thread/database"
	"github.com/Wenfuuu/tpa-web/service-thread/server"
	"github.com/Wenfuuu/tpa-web/service-thread/utils"

	// pb "github.com/Wenfuuu/tpa-web/service-user/proto"
	"fmt"
	"log"
	"net"

	threadpb "github.com/Wenfuuu/tpa-web/pkg/proto/thread"
	"google.golang.org/grpc"
)

func main() {

	db := database.InitDB()
	database.Migrate(db)
	// database.Seed(db)

	// follow service
	followServiceAddr := utils.GetDotEnv("FOLLOW_SERVICE_ADDR")
	followClient, followConn := client.NewFollowClient(followServiceAddr)
	defer followConn.Close()

	port := utils.GetDotEnv("GRPC_PORT")
	if port == "" {
		port = "50054"
	}
	
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	
	grpcServer := grpc.NewServer()
	threadServer := server.NewThreadServer(followClient)
	threadpb.RegisterThreadServiceServer(grpcServer, threadServer)

	log.Printf("Thread service listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}