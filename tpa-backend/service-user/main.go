package main

import (
	"github.com/Wenfuuu/tpa-web/service-user/database"
	"github.com/Wenfuuu/tpa-web/service-user/server"
	"github.com/Wenfuuu/tpa-web/service-user/utils"

	// pb "github.com/Wenfuuu/tpa-web/service-user/proto"
	"fmt"
	"net"

	userpb "github.com/Wenfuuu/tpa-web/pkg/proto/user"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// Initialize logger
	logger := utils.InitLogger()
	defer logger.Sync() // Flush any buffered log entries

	logger.Info("Starting user service")

	db := database.InitDB()
	database.Migrate(db)
	database.Seed(db)

	port := utils.GetDotEnv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err), zap.String("port", port))
	}

	grpcServer := grpc.NewServer()
	userServer := server.NewUserServer()
	userpb.RegisterUserServiceServer(grpcServer, userServer)

	logger.Info("User service started",
		zap.String("port", port),
		zap.String("environment", utils.GetDotEnv("ENVIRONMENT")))

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("Failed to serve", zap.Error(err))
	}
}
