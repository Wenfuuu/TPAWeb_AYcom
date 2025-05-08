package handler

import (
	"context"
	"io"
	"net/http"
	"strings"

	// pb "github.com/Wenfuuu/tpa-web/service-user/proto"
	userpb "github.com/Wenfuuu/tpa-web/pkg/proto/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	Client userpb.UserServiceClient
}

func NewUserHandler(userClient userpb.UserServiceClient) *UserHandler {
	return &UserHandler{
		Client: userClient,
	}
}

// FollowUser godoc
// @Summary Follow a user
// @Description Create a follow relationship between two users
// @Tags Users
// @Accept json
// @Produce json
// @Param request body FollowUserRequest true "Follow User Request"
// @Success 200 {object} userpb.FollowUserResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/follow [post]
func (h *UserHandler) FollowUser(c *gin.Context) {
	var requestBody struct {
		FollowerID         	string `json:"follower_id" binding:"required"`
		FollowingID 		string `json:"following_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.FollowUser(context.Background(), &userpb.FollowUserRequest{
		FollowerId: requestBody.FollowerID,
		FollowingId: requestBody.FollowingID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UnfollowUser godoc
// @Summary Unfollow a user
// @Description Remove a follow relationship between two users
// @Tags Users
// @Produce json
// @Param follower_id query string true "Follower ID"
// @Param following_id query string true "Following ID"
// @Success 200 {string} string
// @Failure 500 {object} map[string]string
// @Router /users/unfollow [delete]
func (h *UserHandler) UnfollowUser(c *gin.Context) {
	followerId := c.Query("follower_id")
	followingId := c.Query("following_id")

	resp, err := h.Client.UnfollowUser(context.Background(), &userpb.UnfollowUserRequest{
		FollowerId: followerId,
		FollowingId: followingId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Message)
}

// GetFollower godoc
// @Summary Get user's followers
// @Description Retrieve a list of users who follow the specified user
// @Tags Users
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {array} userpb.Follow
// @Failure 500 {object} map[string]string
// @Router /users/follower [get]
func (h *UserHandler) GetFollower(c *gin.Context) {
	user_id := c.Query("user_id")

	resp, err := h.Client.GetFollower(context.Background(), &userpb.GetFollowerRequest{
		UserId: user_id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Follows)
}

// GetFollowing godoc
// @Summary Get users followed by a user
// @Description Retrieve a list of users who are followed by the specified user
// @Tags Users
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {array} userpb.Follow
// @Failure 500 {object} map[string]string
// @Router /users/following [get]
func (h *UserHandler) GetFollowing(c *gin.Context) {
	userId := c.Query("user_id")

	resp, err := h.Client.GetFollowing(context.Background(), &userpb.GetFollowingRequest{
		UserId: userId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Follows)
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all registered users
// @Tags Users
// @Produce json
// @Success 200 {array} userpb.User
// @Failure 500 {object} map[string]string
// @Router /users/user [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	resp, err := h.Client.GetAllUsers(context.Background(), &userpb.GetAllUsersRequest{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Users)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Retrieve a specific user by their ID
// @Tags Users
// @Produce json
// @Param id query string true "User ID"
// @Success 200 {object} userpb.User
// @Failure 500 {object} map[string]string
// @Router /users/user-profile [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userId := c.Query("id")

	resp, err := h.Client.GetUser(context.Background(), &userpb.GetUserRequest{
		Id: userId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.User)
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req userpb.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	resp, err := h.Client.ResetPassword(context.Background(), &req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": st.Message()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": resp.Message})
}

func (h *UserHandler) GetSecurityQuestion(c *gin.Context) {
	userId := c.Query("user_id")

	resp, err := h.Client.GetSecurityQuestion(context.Background(), &userpb.GetSecurityQuestionRequest{
		UserId: userId,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			// gRPC error with proper code and message
			c.JSON(http.StatusBadRequest, gin.H{"error": st.Message()})
		} else {
			// Non-gRPC error (unexpected)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Query("email")

	resp, err := h.Client.GetUserByEmail(context.Background(), &userpb.GetUserByEmailRequest{
		Email: email,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.User)
}

// GetWhoToFollow godoc
// @Summary Get users to follow recommendations
// @Description Retrieve a list of users recommended to follow for the specified user
// @Tags Users
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {array} userpb.User
// @Failure 500 {object} map[string]string
// @Router /users/who-to-follow [get]
func (h *UserHandler) GetWhoToFollow(c *gin.Context) {
	userId := c.Query("user_id")

	resp, err := h.Client.GetWhoToFollow(context.Background(), &userpb.GetWhoToFollowRequest{
		// add current logged in later
		UserId: userId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Users)
}

// RegisterPremium godoc
// @Summary Register for premium account
// @Description Register a user for premium account status
// @Tags Users
// @Accept json
// @Produce json
// @Param request body RegisterPremiumRequest true "Premium Registration Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/register-premium [post]
func (h *UserHandler) RegisterPremium(c *gin.Context) {
	var requestBody struct {
		UserID         string `json:"userId" binding:"required"`
		IdentityNumber string `json:"identityNumber" binding:"required"`
		Reason         string `json:"reason" binding:"required"`
		FacePicture    string `json:"facePicture" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.RegisterPremium(context.Background(), &userpb.RegisterPremiumRequest{
		UserId:         requestBody.UserID,
		IdentityNumber: requestBody.IdentityNumber,
		Reason:         requestBody.Reason,
		FacePicture:    requestBody.FacePicture,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": resp.Message,
	})
}

func (h *UserHandler) UploadImage(c *gin.Context) {
	// Get file from request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file type
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File must be an image"})
		return
	}

	// Open the file
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	// Read file contents
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	resp, err := h.Client.UploadImage(context.Background(), &userpb.UploadImageRequest{
		ImageData: fileBytes,
		Filename:  file.Filename,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": resp.ImageUrl,
	})
}

func (h *UserHandler) UploadVideo(c *gin.Context) {
	// Get file from request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file type
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "video/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File must be a video"})
		return
	}

	// Open the file
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	// Read file contents
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	resp, err := h.Client.UploadVideo(context.Background(), &userpb.UploadVideoRequest{
		VideoData: fileBytes,
		Filename:  file.Filename,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": resp.VideoUrl,
	})
}
