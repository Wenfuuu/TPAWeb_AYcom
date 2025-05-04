package handler

import (
	"context"
	"io"
	"net/http"
	"strings"

	// pb "github.com/Wenfuuu/tpa-web/service-user/proto"
	userpb "github.com/Wenfuuu/tpa-web/pkg/proto/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Client userpb.UserServiceClient
}

func NewUserHandler(userClient userpb.UserServiceClient) *UserHandler {
	return &UserHandler{
		Client: userClient,
	}
}

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

func(h *UserHandler) GetAllUsers(c *gin.Context) {
	resp, err := h.Client.GetAllUsers(context.Background(), &userpb.GetAllUsersRequest{})
		
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp.Users)
}

func(h *UserHandler) GetWhoToFollow(c *gin.Context) {
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

func(h *UserHandler) RegisterPremium(c *gin.Context) {
	var requestBody struct {
		UserID   		string 		`json:"userId" binding:"required"`
		IdentityNumber 	string    	`json:"identityNumber" binding:"required"`
		Reason 			string    	`json:"reason" binding:"required"`
		FacePicture		string 		`json:"facePicture" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.RegisterPremium(context.Background(), &userpb.RegisterPremiumRequest{
		UserId: 		requestBody.UserID,
		IdentityNumber: requestBody.IdentityNumber,
		Reason: 		requestBody.Reason,	
		FacePicture: 	requestBody.FacePicture,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": resp.Message,
	})
}

func(h *UserHandler) UploadImage(c *gin.Context) {
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
		Filename: file.Filename,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": resp.ImageUrl,
	})
}