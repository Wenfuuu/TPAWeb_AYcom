package handler

// Swagger request/response types for documentation

// FollowUserRequest represents the request to follow a user
type FollowUserRequest struct {
	FollowerID  string `json:"follower_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	FollowingID string `json:"following_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174001"`
}

// RegisterRequest represents the user registration request
type RegisterRequest struct {
	Name             string `json:"name" binding:"required" example:"John Doe"`
	Username         string `json:"username" binding:"required" example:"johndoe"`
	Email            string `json:"email" binding:"required,email" example:"john@example.com"`
	Password         string `json:"password" binding:"required" example:"password123"`
	Gender           string `json:"gender" binding:"required" example:"male"`
	ProfilePicture   string `json:"profilePicture" binding:"required" example:"https://example.com/profile.jpg"`
	Banner           string `json:"banner" binding:"required" example:"https://example.com/banner.jpg"`
	DateOfBirth      string `json:"dateOfBirth" binding:"required" example:"1990-01-01"`
	SecurityQuestion string `json:"securityQuestion" binding:"required" example:"What was the name of your first pet?"`
	SecurityAnswer   string `json:"securityAnswer" binding:"required" example:"Fluffy"`
	IsSubscribed     bool   `json:"isSubscribed" example:"true"`
}

// LoginRequest represents the user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// RegisterPremiumRequest represents the request to register for premium
type RegisterPremiumRequest struct {
	UserID         string `json:"userId" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	IdentityNumber string `json:"identityNumber" binding:"required" example:"1234567890123456"`
	Reason         string `json:"reason" binding:"required" example:"I want to support the platform"`
	FacePicture    string `json:"facePicture" binding:"required" example:"https://example.com/face.jpg"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}

// MessageResponse represents a success message response
type MessageResponse struct {
	Message string `json:"message" example:"Operation successful"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Operation successful"`
}

// SetUserBanStatusRequest represents the request to set user ban status
type SetUserBanStatusRequest struct {
	UserID string `json:"user_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	Banned bool   `json:"banned" binding:"required" example:"true"`
}

// ResetPasswordRequest represents the request to reset user password
type ResetPasswordRequest struct {
	UserID   string `json:"user_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	Password string `json:"password" binding:"required" example:"newPassword123"`
}

// PremiumRequestID represents the request to accept/reject premium
type PremiumRequestID struct {
	RequestID string `json:"request_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// AcceptReportRequest represents the request to accept a report
type AcceptReportRequest struct {
	ReportID string `json:"report_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	Action   string `json:"action" binding:"required" example:"ban_user"`
}

// RejectReportRequest represents the request to reject a report
type RejectReportRequest struct {
	ReportID string `json:"report_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	Reason   string `json:"reason" binding:"required" example:"Insufficient evidence"`
}

// SendVerificationCodeRequest represents the request to send verification code
type SendVerificationCodeRequest struct {
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
}

// VerifyCodeRequest represents the request to verify code
// VerifyCodeRequest represents the request to verify code
type VerifyCodeRequest struct {
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
	Code  string `json:"code" binding:"required" example:"12345"`
}
