package handler

import (
	"context"
	"net/http"

	authpb "github.com/Wenfuuu/tpa-web/pkg/proto/auth"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Client authpb.AuthServiceClient
}

func NewAuthHandler(authClient authpb.AuthServiceClient) *AuthHandler {
	return &AuthHandler{
		Client: authClient,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var requestBody struct {
		Name     			string `json:"name" binding:"required"`
		Username     		string `json:"username" binding:"required"`
		Email    			string `json:"email" binding:"required,email"`
		Password 			string `json:"password" binding:"required"`
		Gender 				string `json:"gender" binding:"required"`
		ProfilePicture 		string `json:"profilePicture" binding:"required"`
		Banner 				string `json:"banner" binding:"required"`
		DateOfBirth 		string `json:"dateOfBirth" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.Register(context.Background(), &authpb.AuthRegisterRequest{
		Name:           requestBody.Name,
		Username:       requestBody.Username,
		Email:          requestBody.Email,
		Password:       requestBody.Password,
		Gender:         requestBody.Gender,
		ProfilePicture: requestBody.ProfilePicture,
		Banner:         requestBody.Banner,
		DateOfBirth:    requestBody.DateOfBirth,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": resp.User,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var requestBody struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.Login(context.Background(), &authpb.AuthLoginRequest{
		Email:    requestBody.Email,
		Password: requestBody.Password,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"user":          resp.User,
	})
}
