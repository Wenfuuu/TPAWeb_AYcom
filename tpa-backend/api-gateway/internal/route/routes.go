package route

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/Wenfuuu/tpa-web/api-gateway/internal/handler"
	"github.com/Wenfuuu/tpa-web/api-gateway/internal/middleware"

	authpb "github.com/Wenfuuu/tpa-web/pkg/proto/auth"
	chatpb "github.com/Wenfuuu/tpa-web/pkg/proto/chat"
	threadpb "github.com/Wenfuuu/tpa-web/pkg/proto/thread"
	userpb "github.com/Wenfuuu/tpa-web/pkg/proto/user"
	"github.com/gin-gonic/gin"
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

// transcribeAudio forwards audio data to AI service for transcription
func transcribeAudio(c *gin.Context) {
	// Get the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Error reading request body",
		})
		return
	}

	// Create a new request to the AI service
	req, err := http.NewRequest("POST", "http://service-ai:5000/transcribe", bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error creating request to AI service",
		})
		return
	}

	// Copy headers
	req.Header = c.Request.Header

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   "AI service not available",
		})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error reading response from AI service",
		})
		return
	}

	// Set the content type
	c.Header("Content-Type", "application/json")
	c.Status(resp.StatusCode)
	c.Writer.Write(respBody)
}

func InitRoutes(
	r *gin.Engine, 
	userClient userpb.UserServiceClient,
	authClient authpb.AuthServiceClient,
	threadClient threadpb.ThreadServiceClient,
	chatClient chatpb.ChatServiceClient,
) {

	userHandler := handler.NewUserHandler(userClient)
	authHandler := handler.NewAuthHandler(authClient)
	authMiddleware := middleware.NewAuthMiddleware(authClient)
	threadHandler := handler.NewThreadHandler(threadClient)
	chatHandler := handler.NewChatHandler(chatClient, userClient)

	r.POST("/verify-recaptcha", verifyRecaptcha)

	api := r.Group("/api")
	{
		ai := api.Group("/ai")
		{
			ai.POST("/transcribe", transcribeAudio)
		}
		chat := api.Group("/chat")
		{
			chat.GET("/conversations", chatHandler.GetConversations)
			// chat.GET("/conversation", chatHandler.GetConversationById)
			chat.POST("/conversation", chatHandler.CreateConversation)
			chat.GET("/messages", chatHandler.GetMessages)
			chat.POST("/send", chatHandler.SendMessage)
			chat.PUT("/conversation/participants", chatHandler.UpdateParticipant)
			chat.GET("/notifications", chatHandler.GetAllNotifications)
			chat.GET("/mention-notifications", chatHandler.GetMentionNotifications)
		}
		users := api.Group("/users")
		{
			// users.GET("/me", authMiddleware(), getProfile)
			// users.PUT("/:id", authMiddleware(), updateUser)
			// users.DELETE("/:id", authMiddleware(), deleteUser)
			users.GET("/user", userHandler.GetAllUsers)
			users.GET("/user-profile", userHandler.GetUser)
			users.GET("/user-email", userHandler.GetUserByEmail)
			users.GET("/security-question", userHandler.GetSecurityQuestion)
			users.PUT("/reset-password", userHandler.ResetPassword)
			users.POST("/upload/:type", userHandler.UploadImage)
			users.POST("/upload-video/:type", userHandler.UploadVideo)
			users.POST("/follow", userHandler.FollowUser)
			users.DELETE("/unfollow", userHandler.UnfollowUser)
			users.GET("/follower", userHandler.GetFollower)
			users.GET("/following", userHandler.GetFollowing)
			users.GET("/who-to-follow", userHandler.GetWhoToFollow)
			users.POST("/register-premium", userHandler.RegisterPremium)
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
			thread.GET("/thread", threadHandler.GetAllThreads)
			thread.GET("/thread-following", threadHandler.GetFollowingThreads)
			thread.GET("/thread-user", threadHandler.GetUserThreads)
			thread.GET("/thread-user-reply", threadHandler.GetUserReplies)
			thread.GET("/thread-bookmark", threadHandler.GetBookmark)
			thread.GET("/thread-category", threadHandler.GetAllThreadCategories)
			thread.GET("/thread-community", threadHandler.GetCommunityThreads)
			thread.GET("/thread-media", threadHandler.GetThreadMedia)
			thread.GET("/top-thread-community", threadHandler.GetTopCommunityThreads)
			thread.GET("/community", threadHandler.GetAllCommunity)
			thread.GET("/community-member", threadHandler.GetCommunityMembers)
			thread.GET("/community-media", threadHandler.GetCommunityMedia)
			thread.GET("/top-community-member", threadHandler.GetTopCommunityMembers)
			thread.GET("/community-pagination", threadHandler.GetCommunityPagination)
			thread.GET("/community-category", threadHandler.GetAllCommunityCategories)
			thread.GET("/request-join-community", threadHandler.GetJoinCommunity)
			thread.POST("/request-join-community", threadHandler.RequestJoinCommunity)
			thread.GET("/request-create-community", threadHandler.GetCreateCommunity)
			thread.POST("/request-create-community", threadHandler.RequestCreateCommunity)
			thread.POST("/accept-join-community", threadHandler.AcceptJoinCommunity)
			thread.POST("/reject-join-community", threadHandler.RejectJoinCommunity)
			thread.GET("/pending-community", threadHandler.GetPendingCommunity)
			thread.GET("/joined-community", threadHandler.GetJoinedCommunity)
			thread.GET("/thread-file", threadHandler.GetAllThreadFile)
			thread.POST("/post", threadHandler.PostThread)
			thread.POST("/reply", threadHandler.ReplyThread)
			thread.GET("/repost", threadHandler.GetRepostThreads)
			thread.POST("/repost", threadHandler.RepostThread)
			thread.POST("/unrepost", threadHandler.UnrepostThread)
			thread.GET("/like", threadHandler.GetLikedThreads)
			thread.POST("/like", threadHandler.LikeThread)
			thread.POST("/unlike", threadHandler.UnlikeThread)
			thread.GET("/bookmark", threadHandler.GetBookmarkThreads)
			thread.POST("/bookmark", threadHandler.BookmarkThread)
			thread.POST("/unbookmark", threadHandler.UnbookmarkThread)
			thread.GET("/thread-detail", threadHandler.GetThreadById)
			thread.GET("/community-detail", threadHandler.GetCommunityById)
			thread.GET("/thread-reply", threadHandler.GetReplies)
			thread.GET("/top-hashtag", threadHandler.GetTopHashtags)
		}
	}

}