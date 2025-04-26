package handler

import (
	"context"
	"net/http"

	pb "github.com/Wenfuuu/tpa-web/service-user/proto"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Client pb.UserServiceClient
}

func NewUserHandler(userClient pb.UserServiceClient) *UserHandler {
	return &UserHandler{
		Client: userClient,
	}
}

func(h *UserHandler) GetAllUsers(c *gin.Context) {
	resp, err := h.Client.GetAllUsers(context.Background(), &pb.GetAllUsersRequest{})
		
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp.Users)
}

func(h* UserHandler) LoginUser(c *gin.Context) {
	var requestBody struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.Login(context.Background(), &pb.LoginRequest{
		Email: requestBody.Email,
		Password: requestBody.Password,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": resp.Token,
		"user":  resp.User,
	})
}

func(h *UserHandler) RegisterUser(c *gin.Context) {
	var requestBody struct {
		Name     			string `json:"name" binding:"required"`
		Username     		string `json:"username" binding:"required"`
		Email    			string `json:"email" binding:"required,email"`
		Password 			string `json:"password" binding:"required"`
		Gender 				string `json:"gender" binding:"required"`
		ProfilePicture 		string `json:"profilePicture" binding:"required"`
		Banner 				string `json:"banner" binding:"required"`
		DateOfBirth 		string `json:"dateOfBirth" binding:"required"`
		IsActive        	bool   `json:"isActive"` // Optional, default to true in backend if needed
		IsBanned        	bool   `json:"isBanned"` // Optional, default to false in backend
		Role            	string `json:"role"`     // Optional, default to "user"
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.Register(context.Background(), &pb.RegisterRequest{
		Name:           requestBody.Name,
		Username:       requestBody.Username,
		Email:          requestBody.Email,
		Password:       requestBody.Password,
		Gender:         requestBody.Gender,
		ProfilePicture: requestBody.ProfilePicture,
		Banner:         requestBody.Banner,
		DateOfBirth:    requestBody.DateOfBirth,
		IsActive:       requestBody.IsActive,
		IsBanned:       requestBody.IsBanned,
		Role:           requestBody.Role,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  resp.User,
	})
}