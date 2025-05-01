package main

import (
	"github.com/Wenfuuu/tpa-web/service-thread/database"
	"github.com/Wenfuuu/tpa-web/service-thread/utils"
	"github.com/Wenfuuu/tpa-web/service-thread/server"
	// pb "github.com/Wenfuuu/tpa-web/service-user/proto"
	threadpb "github.com/Wenfuuu/tpa-web/pkg/proto/thread"
	"google.golang.org/grpc"
	"fmt"
	"log"
	"net"
)

func main() {

	db := database.InitDB()
	database.Migrate(db)
	// database.Seed(db)

	port := utils.GetDotEnv("GRPC_PORT")
	if port == "" {
		port = "50054"
	}
	
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	
	grpcServer := grpc.NewServer()
	threadServer := server.NewThreadServer()
	threadpb.RegisterThreadServiceServer(grpcServer, threadServer)

	log.Printf("Thread service listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}