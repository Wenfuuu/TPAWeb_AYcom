package main

import (
	"github.com/Wenfuuu/tpa-web/service-user/database"
	"github.com/Wenfuuu/tpa-web/service-user/utils"
	"github.com/Wenfuuu/tpa-web/service-user/server"
	pb "github.com/Wenfuuu/tpa-web/service-user/proto"
	"google.golang.org/grpc"
	"fmt"
	"log"
	"net"
)

func main() {

	db := database.InitDB()
	database.Migrate(db)
	database.Seed(db)

	port := utils.GetDotEnv("SERVICE_USER_PORT")
	if port == "" {
		port = "50051"
	}
	
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	
	grpcServer := grpc.NewServer()
	userServer := server.NewUserServer()
	pb.RegisterUserServiceServer(grpcServer, userServer)

	log.Printf("User service listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}