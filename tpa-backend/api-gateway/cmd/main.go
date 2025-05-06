package main

import (
	"log"
	"time"

	"github.com/Wenfuuu/tpa-web/api-gateway/internal/client"
	"github.com/Wenfuuu/tpa-web/api-gateway/internal/route"
	"github.com/Wenfuuu/tpa-web/api-gateway/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// user service
	userServiceAddr := utils.GetDotEnv("USER_SERVICE_ADDR")
	userClient, userConn := client.NewUserClient(userServiceAddr)
	defer userConn.Close()

	// auth service
	authServiceAddr := utils.GetDotEnv("AUTH_SERVICE_ADDR")
	authClient, authConn := client.NewAuthClient(authServiceAddr)
	defer authConn.Close()

	// thread service
	threadServiceAddr := utils.GetDotEnv("THREAD_SERVICE_ADDR")
	threadClient, threadConn := client.NewThreadClient(threadServiceAddr)
	defer threadConn.Close()

	// chat service
	chatServiceAddr := utils.GetDotEnv("CHAT_SERVICE_ADDR")
	chatClient, chatConn := client.NewChatClient(chatServiceAddr)
	defer chatConn.Close()
	
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Refresh-token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//setup route
	route.InitRoutes(r, userClient, authClient, threadClient, chatClient)

	port := utils.GetDotEnv("PORT")
	log.Printf("API Gateway listening on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}