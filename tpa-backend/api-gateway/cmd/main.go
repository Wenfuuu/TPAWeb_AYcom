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
	
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Vite default port
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//setup route
	route.InitRoutes(r, userClient, authClient)

	port := utils.GetDotEnv("PORT")
	log.Printf("API Gateway listening on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}