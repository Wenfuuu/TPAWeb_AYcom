package route

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/Wenfuuu/tpa-web/api-gateway/internal/handler"
	"github.com/Wenfuuu/tpa-web/api-gateway/internal/middleware"

	// pb "github.com/Wenfuuu/tpa-web/service-user/proto"
	authpb "github.com/Wenfuuu/tpa-web/pkg/proto/auth"
	userpb "github.com/Wenfuuu/tpa-web/pkg/proto/user"
	threadpb "github.com/Wenfuuu/tpa-web/pkg/proto/thread"
	followpb "github.com/Wenfuuu/tpa-web/pkg/proto/follow"
	"github.com/gin-gonic/gin"
	// "github.com/mojocn/base64Captcha"
)

const recaptchaSecret = "6Le5ySYrAAAAAKMMRBku5GFQxZnzfSjftwKO9yQS"

type RecaptchaResponse struct {
    Success     bool     `json:"success"`
    ChallengeTS string   `json:"challenge_ts"`
    Hostname    string   `json:"hostname"`
    ErrorCodes  []string `json:"error-codes,omitempty"`
}

func verifyRecaptcha(c *gin.Context) {
	// Define a struct for expected JSON request body
    type RequestBody struct {
        Token string `json:"token"`
    }

    var reqBody RequestBody
    if err := c.ShouldBindJSON(&reqBody); err != nil || reqBody.Token == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // Post to Google's reCAPTCHA verification API using application/x-www-form-urlencoded
    resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
        url.Values{
            "secret":   {recaptchaSecret},
            "response": {reqBody.Token},
        })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify reCAPTCHA"})
        return
    }
    defer resp.Body.Close()

    var recaptchaResp RecaptchaResponse
    if err := json.NewDecoder(resp.Body).Decode(&recaptchaResp); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode reCAPTCHA response"})
        return
    }

    c.JSON(http.StatusOK, recaptchaResp)
}

func InitRoutes(
	r *gin.Engine, 
	userClient userpb.UserServiceClient,
	authClient authpb.AuthServiceClient,
	threadClient threadpb.ThreadServiceClient,
	followClient followpb.FollowServiceClient,
) {

	userHandler := handler.NewUserHandler(userClient)
	authHandler := handler.NewAuthHandler(authClient)
	authMiddleware := middleware.NewAuthMiddleware(authClient)
	threadHandler := handler.NewThreadHandler(threadClient)
	followHandler := handler.NewFollowHandler(followClient)

	r.POST("/verify-recaptcha", verifyRecaptcha)

	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			// users.GET("/me", authMiddleware(), getProfile)
			// users.PUT("/:id", authMiddleware(), updateUser)
			// users.DELETE("/:id", authMiddleware(), deleteUser)
			users.GET("/user", userHandler.GetAllUsers) // NEED GUARD
			users.POST("/upload/:type", userHandler.UploadImage) // NEED GUARD
		}
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/guest", authMiddleware.CheckGuest())
			auth.GET("/check", authMiddleware.CheckAuth())
			auth.GET("/refresh", authHandler.RefreshToken)
			auth.GET("/get-auth-user", authHandler.GetAuthUser)
			auth.GET("/admin", authMiddleware.CheckAdmin())
			auth.POST("/send-code", authHandler.SendVerificationCode)
			auth.POST("/verify-code", authHandler.VerifyCode)
		}
		thread := api.Group("/threads")
		{
			thread.GET("/thread", threadHandler.GetAllThreads) // NEED GUARD
			thread.GET("/thread-following", threadHandler.GetFollowingThreads) // NEED GUARD
			thread.GET("/thread-category", threadHandler.GetAllThreadCategories) // NEED GUARD
			thread.GET("/community", threadHandler.GetAllCommunity) // NEED GUARD
			thread.GET("/thread-file", threadHandler.GetAllThreadFile) // NEED GUARD
			thread.POST("/post", threadHandler.PostThread)
			thread.POST("/like", threadHandler.LikeThread)
			thread.POST("/unlike", threadHandler.UnlikeThread)
		}
		follow := api.Group("/follows")
		{
			follow.GET("/follower", followHandler.GetFollower)
			follow.GET("/following", followHandler.GetFollowing)
		}
	}

}