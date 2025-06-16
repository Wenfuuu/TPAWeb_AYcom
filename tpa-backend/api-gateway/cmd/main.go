package main

import (
	"log"
	"time"

	"github.com/Wenfuuu/tpa-web/api-gateway/internal/client"
	"github.com/Wenfuuu/tpa-web/api-gateway/internal/route"
	"github.com/Wenfuuu/tpa-web/api-gateway/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// Swagger
	docs "github.com/Wenfuuu/tpa-web/api-gateway/cmd/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title TPA Web API Gateway
// @version 1.0
// @description This is the API Gateway for the TPA Web microservices project.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Initialize Swagger docs
	docs.SwaggerInfo.BasePath = "/api"

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

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//setup route
	route.InitRoutes(r, userClient, authClient, threadClient, chatClient)

	port := utils.GetDotEnv("PORT")
	log.Printf("API Gateway listening on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
