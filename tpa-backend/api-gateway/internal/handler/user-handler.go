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
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/follow [post]
func (h *UserHandler) FollowUser(c *gin.Context) {
	var requestBody struct {
		FollowerID  string `json:"follower_id" binding:"required"`
		FollowingID string `json:"following_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.FollowUser(context.Background(), &userpb.FollowUserRequest{
		FollowerId:  requestBody.FollowerID,
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
// @Success 200 {object} MessageResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/unfollow [delete]
func (h *UserHandler) UnfollowUser(c *gin.Context) {
	followerId := c.Query("follower_id")
	followingId := c.Query("following_id")

	resp, err := h.Client.UnfollowUser(context.Background(), &userpb.UnfollowUserRequest{
		FollowerId:  followerId,
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
// @Success 200 {array} object
// @Failure 500 {object} ErrorResponse
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
// @Success 200 {array} object
// @Failure 500 {object} ErrorResponse
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
// @Success 200 {array} object
// @Failure 500 {object} ErrorResponse
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
// @Success 200 {object} object
// @Failure 500 {object} ErrorResponse
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

// SetUserBanStatus godoc
// @Summary Set user ban status
// @Description Update the ban status of a user (admin function)
// @Tags Users
// @Accept json
// @Produce json
// @Param request body SetUserBanStatusRequest true "Ban Status Request"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/ban [put]
func (h *UserHandler) SetUserBanStatus(c *gin.Context) {
	var req userpb.SetUserBanStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	resp, err := h.Client.SetUserBanStatus(context.Background(), &req)
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

// ResetPassword godoc
// @Summary Reset user password
// @Description Reset a user's password using security question
// @Tags Users
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Reset Password Request"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/reset-password [put]
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

// GetSecurityQuestion godoc
// @Summary Get user's security question
// @Description Retrieve the security question for a user by user ID
// @Tags Users
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {object} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/security-question [get]
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

// GetUserByEmail godoc
// @Summary Get user by email
// @Description Retrieve a user by their email address
// @Tags Users
// @Produce json
// @Param email query string true "User Email"
// @Success 200 {object} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/user-email [get]
func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Query("email")

	resp, err := h.Client.GetUserByEmail(context.Background(), &userpb.GetUserByEmailRequest{
		Email: email,
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

	c.JSON(http.StatusOK, resp.User)
}

// GetWhoToFollow godoc
// @Summary Get users to follow recommendations
// @Description Retrieve a list of users recommended to follow for the specified user
// @Tags Users
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {array} object
// @Failure 500 {object} ErrorResponse
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

// AcceptPremium godoc
// @Summary Accept premium account request
// @Description Accept a user's premium account request (admin function)
// @Tags Users
// @Accept json
// @Produce json
// @Param request body PremiumRequestID true "Premium Request ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/accept-premium [post]
func (h *UserHandler) AcceptPremium(c *gin.Context) {
	var req userpb.PremiumRequestId

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	resp, err := h.Client.AcceptPremiumRequest(context.Background(), &req)
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

// RejectPremium godoc
// @Summary Reject premium account request
// @Description Reject a user's premium account request (admin function)
// @Tags Users
// @Accept json
// @Produce json
// @Param request body PremiumRequestID true "Premium Request ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/reject-premium [post]
func (h *UserHandler) RejectPremium(c *gin.Context) {
	var req userpb.PremiumRequestId

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	resp, err := h.Client.RejectPremiumRequest(context.Background(), &req)
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

// GetPremiumRequests godoc
// @Summary Get all premium account requests
// @Description Retrieve all pending premium account requests (admin function)
// @Tags Users
// @Produce json
// @Success 200 {array} object
// @Failure 500 {object} ErrorResponse
// @Router /users/premium-request [get]
func (h *UserHandler) GetPremiumRequests(c *gin.Context) {
	resp, err := h.Client.GetPremiumRequests(context.Background(), &userpb.GetPremiumRequestsRequest{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.PremiumRequests)
}

// ReportUser godoc
// @Summary Report a user
// @Description Report a user for inappropriate behavior
// @Tags Users
// @Accept json
// @Produce json
// @Param request body ReportUserRequest true "Report User Request"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/report [post]
func (h *UserHandler) ReportUser(c *gin.Context) {
	var requestBody struct {
		UserID     string `json:"user_id" binding:"required"`
		ReportedID string `json:"reported_id" binding:"required"`
		Reason     string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.ReportUser(context.Background(), &userpb.ReportUserRequest{
		UserId:     requestBody.UserID,
		ReportedId: requestBody.ReportedID,
		Reason:     requestBody.Reason,
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

	c.JSON(http.StatusOK, gin.H{"message": resp.Message})
}

// AcceptReport godoc
// @Summary Accept a user report
// @Description Accept a user report and take action (admin function)
// @Tags Users
// @Accept json
// @Produce json
// @Param request body AcceptReportRequest true "Accept Report Request"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/accept-report [post]
func (h *UserHandler) AcceptReport(c *gin.Context) {
	var req userpb.AcceptReportRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	resp, err := h.Client.AcceptReport(context.Background(), &req)
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

// RejectReport godoc
// @Summary Reject a user report
// @Description Reject a user report (admin function)
// @Tags Users
// @Accept json
// @Produce json
// @Param request body RejectReportRequest true "Reject Report Request"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/reject-report [post]
func (h *UserHandler) RejectReport(c *gin.Context) {
	var req userpb.RejectReportRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	resp, err := h.Client.RejectReport(context.Background(), &req)
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

// GetReports godoc
// @Summary Get all user reports
// @Description Retrieve all user reports (admin function)
// @Tags Users
// @Produce json
// @Success 200 {array} object
// @Failure 500 {object} ErrorResponse
// @Router /users/reports [get]
func (h *UserHandler) GetReports(c *gin.Context) {
	resp, err := h.Client.GetReports(context.Background(), &userpb.GetReportsRequest{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Reports)
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
		st, ok := status.FromError(err)
		if ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": st.Message()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": resp.Message,
	})
}

// UploadImage godoc
// @Summary Upload an image file
// @Description Upload an image file and get the URL back
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Param type path string true "Upload type" example:"profile"
// @Param file formData file true "Image file to upload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/upload/{type} [post]
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

// UploadVideo godoc
// @Summary Upload a video file
// @Description Upload a video file and get the URL back
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Param type path string true "Upload type" example:"content"
// @Param file formData file true "Video file to upload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/upload-video/{type} [post]
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
