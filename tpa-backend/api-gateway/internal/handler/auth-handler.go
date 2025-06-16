package handler

import (
	"context"
	"net/http"

	authpb "github.com/Wenfuuu/tpa-web/pkg/proto/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	Client authpb.AuthServiceClient
}

func NewAuthHandler(authClient authpb.AuthServiceClient) *AuthHandler {
	return &AuthHandler{
		Client: authClient,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration data"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var requestBody struct {
		Name             string `json:"name" binding:"required"`
		Username         string `json:"username" binding:"required"`
		Email            string `json:"email" binding:"required,email"`
		Password         string `json:"password" binding:"required"`
		Gender           string `json:"gender" binding:"required"`
		ProfilePicture   string `json:"profilePicture" binding:"required"`
		Banner           string `json:"banner" binding:"required"`
		DateOfBirth      string `json:"dateOfBirth" binding:"required"`
		SecurityQuestion string `json:"securityQuestion" binding:"required"`
		SecurityAnswer   string `json:"securityAnswer" binding:"required"`
		IsSubscribed     bool   `json:"isSubscribed"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.Register(context.Background(), &authpb.AuthRegisterRequest{
		Name:             requestBody.Name,
		Username:         requestBody.Username,
		Email:            requestBody.Email,
		Password:         requestBody.Password,
		Gender:           requestBody.Gender,
		ProfilePicture:   requestBody.ProfilePicture,
		Banner:           requestBody.Banner,
		DateOfBirth:      requestBody.DateOfBirth,
		SecurityQuestion: requestBody.SecurityQuestion,
		SecurityAnswer:   requestBody.SecurityAnswer,
		IsSubscribed:     requestBody.IsSubscribed,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": st.Message()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": resp.User,
	})
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return access tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
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
		st, ok := status.FromError(err)
		if ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": st.Message()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"user":          resp.User,
	})
}

// GetAuthUser godoc
// @Summary Get authenticated user
// @Description Get current authenticated user information
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Router /auth/get-auth-user [get]
func (h *AuthHandler) GetAuthUser(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header required"})
		c.Abort()
		return
	}

	accessTokenString := authHeader[len("Bearer "):]
	resp, err := h.Client.GetAuthUser(context.Background(), &authpb.GetAuthUserRequest{
		AccessToken: accessTokenString,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": resp.User,
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get new access and refresh tokens using refresh token
// @Tags Authentication
// @Produce json
// @Param Refresh-token header string true "Refresh Token"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Router /auth/refresh [get]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("Refresh-token")
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Refresh-token header required"})
		c.Abort()
		return
	}

	resp, err := h.Client.RefreshToken(context.Background(), &authpb.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"user":          resp.User,
	})
}

// SendVerificationCode godoc
// @Summary Send verification code
// @Description Send verification code to user's email for registration
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body SendVerificationCodeRequest true "Email for verification"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/send-code [post]
func (h *AuthHandler) SendVerificationCode(c *gin.Context) {
	var requestBody struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.SendVerificationCode(context.Background(), &authpb.SendVerificationCodeRequest{
		Email: requestBody.Email,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"send_code": resp.Success,
	})
}

// VerifyCode godoc
// @Summary Verify verification code
// @Description Verify the verification code sent to user's email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body VerifyCodeRequest true "Email and verification code"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/verify-code [post]
func (h *AuthHandler) VerifyCode(c *gin.Context) {
	var requestBody struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.VerifyCode(context.Background(), &authpb.VerifyCodeRequest{
		Email: requestBody.Email,
		Code:  requestBody.Code,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Code is not valid or already expired"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"verified": resp.Verified,
	})
}
