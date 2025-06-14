// project/backend/service-user/handler/user_handler.go
package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Wenfuuu/tpa-web/service-user/model"
	"github.com/Wenfuuu/tpa-web/service-user/utils"
	"github.com/google/uuid"

	// pb "github.com/Wenfuuu/tpa-web/service-user/proto"
	userpb "github.com/Wenfuuu/tpa-web/pkg/proto/user"
	"github.com/Wenfuuu/tpa-web/service-user/repository"

	// "github.com/Wenfuuu/tpa-web/service-user/utils"
	storage_go "github.com/supabase-community/storage-go"
	supa "github.com/supabase-community/supabase-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userRepo: repository.NewUserRepository(),
	}
}

func (h *UserHandler) FollowUser(ctx context.Context, req *userpb.FollowUserRequest) (*userpb.Follow, error) {
	followerID, err := uuid.Parse(req.FollowerId)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	followingID, err := uuid.Parse(req.FollowingId)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	follow := model.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	if err := h.userRepo.FollowUser(&follow); err != nil {
		return nil, status.Error(codes.Internal, "Failed to follow user")
	}

	return h.followToProto(&follow), nil
}

func (h *UserHandler) UnfollowUser(ctx context.Context, req *userpb.UnfollowUserRequest) (*userpb.UnfollowUserResponse, error) {
	err := h.userRepo.UnfollowUser(req.FollowerId, req.FollowingId)

	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}

	return &userpb.UnfollowUserResponse{
		Message: "Success unfollow user!",
	}, nil
}

func (h *UserHandler) GetFollower(ctx context.Context, req *userpb.GetFollowerRequest) (*userpb.FollowResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "User ID is required")
	}

	modelFollows, err := h.userRepo.GetFollower(req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}

	var pbFollows []*userpb.Follow
	for _, f := range modelFollows {
		pbFollows = append(pbFollows, h.followToProto(&f))
	}

	return &userpb.FollowResponse{
		Follows: pbFollows,
	}, nil
}

func (h *UserHandler) GetFollowing(ctx context.Context, req *userpb.GetFollowingRequest) (*userpb.FollowResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "User ID is required")
	}

	modelFollows, err := h.userRepo.GetFollowing(req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}

	var pbFollows []*userpb.Follow
	for _, f := range modelFollows {
		pbFollows = append(pbFollows, h.followToProto(&f))
	}

	return &userpb.FollowResponse{
		Follows: pbFollows,
	}, nil
}

func (h *UserHandler) followToProto(follow *model.Follow) *userpb.Follow {
	return &userpb.Follow{
		FollowerId:  follow.FollowerID.String(),
		FollowingId: follow.FollowingID.String(),
	}
}

func (h *UserHandler) AcceptPremium(ctx context.Context, req *userpb.PremiumRequestId) (*userpb.ActionResponse, error) {
	if err := h.userRepo.AcceptPremium(req.RequestId); err != nil {
		return nil, status.Error(codes.Internal, "Failed to accept community creation request")
	}

	return &userpb.ActionResponse{
		Message: "Premium request accepted successfully",
	}, nil
}

func (h *UserHandler) RejectPremium(ctx context.Context, req *userpb.PremiumRequestId) (*userpb.ActionResponse, error) {
	if err := h.userRepo.RejectPremium(req.RequestId); err != nil {
		return nil, status.Error(codes.Internal, "Failed to reject community creation request")
	}

	return &userpb.ActionResponse{
		Message: "Premium request rejected successfully",
	}, nil
}

func (h *UserHandler) GetPremiumRequests(ctx context.Context, req *userpb.GetPremiumRequestsRequest) (*userpb.GetPremiumRequestsResponse, error) {
	modelRequests, err := h.userRepo.GetPremiumRequests()
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}

	var pbRequests []*userpb.PremiumRequest
	for _, r := range modelRequests {
		pbRequests = append(pbRequests, h.premiumRequestToProto(&r))
	}

	return &userpb.GetPremiumRequestsResponse{
		PremiumRequests: pbRequests,
	}, nil
}

func (h *UserHandler) AcceptReport(ctx context.Context, req *userpb.AcceptReportRequest) (*userpb.ActionResponse, error) {
	if err := h.userRepo.AcceptReport(req.ReportId); err != nil {
		return nil, status.Error(codes.Internal, "Failed to accept report")
	}

	return &userpb.ActionResponse{
		Message: "Report accepted successfully",
	}, nil
}

func (h *UserHandler) RejectReport(ctx context.Context, req *userpb.RejectReportRequest) (*userpb.ActionResponse, error) {
	if err := h.userRepo.RejectReport(req.ReportId); err != nil {
		return nil, status.Error(codes.Internal, "Failed to reject report")
	}

	return &userpb.ActionResponse{
		Message: "Report rejected successfully",
	}, nil
}

func (h *UserHandler) GetReports(ctx context.Context, req *userpb.GetReportsRequest) (*userpb.GetReportsResponse, error) {
	modelRequests, err := h.userRepo.GetReports()
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}

	var pbReports []*userpb.Report
	for _, r := range modelRequests {
		pbReports = append(pbReports, h.reportToProto(&r))
	}

	return &userpb.GetReportsResponse{
		Reports: pbReports,
	}, nil
}

func (h *UserHandler) RegisterPremium(ctx context.Context, req *userpb.RegisterPremiumRequest) (*userpb.RegisterPremiumResponse, error) {
	// hash id
	hashedId, err := utils.EncryptPassword(req.IdentityNumber)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to hash password")
	}

	// check already registered
	exists, err := h.userRepo.AlreadyRegisteredPremium(req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}
	if exists {
		return nil, status.Error(codes.AlreadyExists, "User already requested a premium request")
	}

	// check already premium
	exists, err = h.userRepo.AlreadyPremium(req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}
	if exists {
		return nil, status.Error(codes.AlreadyExists, "User already premium")
	}

	request := model.PremiumRequest{
		UserID:         uuid.MustParse(req.UserId),
		IdentityNumber: hashedId,
		Reason:         req.Reason,
		FacePicture:    req.FacePicture,
	}

	if err := h.userRepo.RegisterPremium(&request); err != nil {
		return nil, status.Error(codes.Internal, "Failed to register premium")
	}

	return &userpb.RegisterPremiumResponse{
		Message: "Success request premium",
	}, nil
}

func (h *UserHandler) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.UserResponse, error) {
	// Validate input
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// Check if username already exists
	exists, err := h.userRepo.UsernameExists(req.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}
	if exists {
		return nil, status.Error(codes.AlreadyExists, "Username already registered")
	}

	// Check if email already exists
	exists, err = h.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}
	if exists {
		return nil, status.Error(codes.AlreadyExists, "Email already registered")
	}

	// log.Println("Received password:", req.Password)

	// Create new user
	user := model.User{
		Name:           req.Name,
		Username:       req.Username,
		Email:          req.Email,
		Password:       req.Password,
		Gender:         req.Gender,
		ProfilePicture: req.ProfilePicture,
		Banner:         req.Banner,
		DateOfBirth:    req.DateOfBirth,
		IsActive:       req.IsActive,
		IsBanned:       req.IsBanned,
		Role:           req.Role,
		IsSubscribed:   req.IsSubscribed,
		IsPremium:      req.IsPremium,
		Bio:            req.Bio,
	}

	if err := h.userRepo.Create(&user); err != nil {
		return nil, status.Error(codes.Internal, "Failed to create user")
	}

	// log.Printf(req.SecurityQuestion)
	securityQuestion := model.SecurityQuestion{
		UserID:   user.ID,
		Question: req.SecurityQuestion,
		Answer:   req.SecurityAnswer,
	}

	if err := h.userRepo.CreateSecurityQuestion(&securityQuestion); err != nil {
		return nil, status.Error(codes.Internal, "Failed to save security question")
	}

	return h.userToResponse(&user), nil
}

func (h *UserHandler) GetUserByEmail(ctx context.Context, req *userpb.GetUserByEmailRequest) (*userpb.UserResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "Email is required")
	}

	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, "Database error")
	}

	if !user.IsActive {
		return nil, status.Error(codes.PermissionDenied, "User account is inactive")
	}

	if user.IsBanned {
		return nil, status.Error(codes.PermissionDenied, "User account is banned")
	}

	return h.userToResponse(user), nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "User ID is required")
	}

	user, err := h.userRepo.FindByID(req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, "Database error")
	}

	return h.userToResponse(user), nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "User ID is required")
	}

	user, err := h.userRepo.FindByID(req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, "Database error")
	}

	// Update fields if provided
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := h.userRepo.Update(user); err != nil {
		return nil, status.Error(codes.Internal, "Failed to update user")
	}

	return h.userToResponse(user), nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "User ID is required")
	}

	rowsAffected, err := h.userRepo.Delete(req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}
	if rowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	return &userpb.DeleteUserResponse{
		Success: true,
	}, nil
}

func (h *UserHandler) SetUserBanStatus(ctx context.Context, req *userpb.SetUserBanStatusRequest) (*userpb.SetUserBanStatusResponse, error) {
	user, err := h.userRepo.FindByID(req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, "Database error")
	}

	user.IsBanned = req.Banned

	if err := h.userRepo.Update(user); err != nil {
		return nil, status.Error(codes.Internal, "Failed to update password")
	}

	return &userpb.SetUserBanStatusResponse{
		Message: "Success update user ban status",
	}, nil
}

func (h *UserHandler) ResetPassword(ctx context.Context, req *userpb.ResetPasswordRequest) (*userpb.ResetPasswordResponse, error) {
	user, err := h.userRepo.FindByID(req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, "Database error")
	}

	// Compare new password with old one
	if utils.ComparePassword(user.Password, req.Password) {
		return nil, status.Error(codes.InvalidArgument, "New password cannot be the same as the old password")
	}

	// Hash new password
	hashedPassword, err := utils.EncryptPassword(req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to hash password")
	}

	// Update in DB
	user.Password = hashedPassword
	if err := h.userRepo.Update(user); err != nil {
		return nil, status.Error(codes.Internal, "Failed to update password")
	}

	return &userpb.ResetPasswordResponse{Message: "Password updated successfully"}, nil
}

func (h *UserHandler) GetSecurityQuestion(ctx context.Context, req *userpb.GetSecurityQuestionRequest) (*userpb.SecurityQuestionResponse, error) {
	question, err := h.userRepo.GetSecurityQuestion(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "security question not found: %v", err)
	}

	return &userpb.SecurityQuestionResponse{
		UserId:           question.UserID.String(),
		SecurityQuestion: question.Question,
		SecurityAnswer:   question.Answer,
	}, nil
}

func (h *UserHandler) GetWhoToFollow(ctx context.Context, req *userpb.GetWhoToFollowRequest) (*userpb.GetAllUsersResponse, error) {
	modelUsers, err := h.userRepo.GetWhoToFollow(req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}

	var pbUsers []*userpb.User
	for _, u := range modelUsers {
		pbUsers = append(pbUsers, h.userToProto(&u))
	}

	return &userpb.GetAllUsersResponse{
		Users: pbUsers,
	}, nil
}

func (h *UserHandler) GetAllUsers(ctx context.Context, req *userpb.GetAllUsersRequest) (*userpb.GetAllUsersResponse, error) {
	modelUsers, err := h.userRepo.GetAllUsers()
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}

	var pbUsers []*userpb.User
	for _, u := range modelUsers {
		pbUsers = append(pbUsers, h.userToProto(&u))
	}

	return &userpb.GetAllUsersResponse{
		Users: pbUsers,
	}, nil
}

func (h *UserHandler) UploadImage(ctx context.Context, req *userpb.UploadImageRequest) (*userpb.UploadImageResponse, error) {
	supabaseURL := utils.GetDotEnv("SUPABASE_URL")
	supabaseKey := utils.GetDotEnv("SUPABASE_KEY")
	bucketName := utils.GetDotEnv("BUCKET_NAME")

	supabase, err := supa.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to initialize supabase client")
	}

	// Generate unique filename with timestamp
	fileExt := strings.ToLower(filepath.Ext(req.Filename))
	fileBase := strings.TrimSuffix(req.Filename, fileExt)
	timestamp := time.Now().UnixNano()
	uniqueFilename := fmt.Sprintf("%s_%d%s", fileBase, timestamp, fileExt)

	contentType := "application/octet-stream" // default
	switch fileExt {
	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	case ".svg":
		contentType = "image/svg+xml"
	}

	// Upload to Supabase Storage
	fileReader := bytes.NewReader(req.ImageData)
	_, err = supabase.Storage.UploadFile(bucketName, uniqueFilename, fileReader, storage_go.FileOptions{
		CacheControl: stringPtr("3600"),
		ContentType:  &contentType,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to upload image")
	}

	// Get public URL
	publicURL := supabase.Storage.GetPublicUrl(bucketName, uniqueFilename).SignedURL

	return &userpb.UploadImageResponse{
		ImageUrl: publicURL,
	}, nil
}

func (h *UserHandler) UploadVideo(ctx context.Context, req *userpb.UploadVideoRequest) (*userpb.UploadVideoResponse, error) {
	supabaseURL := utils.GetDotEnv("SUPABASE_URL")
	supabaseKey := utils.GetDotEnv("SUPABASE_KEY")
	bucketName := utils.GetDotEnv("BUCKET_NAME")

	supabase, err := supa.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to initialize supabase client")
	}

	// Generate unique filename with timestamp
	fileExt := strings.ToLower(filepath.Ext(req.Filename))
	fileBase := strings.TrimSuffix(req.Filename, fileExt)
	timestamp := time.Now().UnixNano()
	uniqueFilename := fmt.Sprintf("%s_%d%s", fileBase, timestamp, fileExt)

	contentType := "application/octet-stream" // default
	switch fileExt {
	case ".mp4":
		contentType = "video/mp4"
	case ".avi":
		contentType = "video/x-msvideo"
	case ".mov":
		contentType = "video/quicktime"
	case ".mkv":
		contentType = "video/x-matroska"
	default:
		return nil, status.Error(codes.InvalidArgument, "Unsupported video format")
	}

	// Upload to Supabase Storage
	fileReader := bytes.NewReader(req.VideoData)
	_, err = supabase.Storage.UploadFile(bucketName, uniqueFilename, fileReader, storage_go.FileOptions{
		CacheControl: stringPtr("3600"),
		ContentType:  &contentType,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to upload video")
	}

	// Get public URL
	publicURL := supabase.Storage.GetPublicUrl(bucketName, uniqueFilename).SignedURL

	return &userpb.UploadVideoResponse{
		VideoUrl: publicURL,
	}, nil
}

func stringPtr(s string) *string {
	return &s
}

// Helper functions
func (h *UserHandler) reportToProto(req *model.Report) *userpb.Report {
	return &userpb.Report{
		Id:         req.ID.String(),
		UserId:     req.UserID.String(),
		ReportedId: req.ReportedID.String(),
		Reason:     req.Reason,
		CreatedAt:  req.CreatedAt.Format(time.RFC3339),
	}
}

func (h *UserHandler) premiumRequestToProto(req *model.PremiumRequest) *userpb.PremiumRequest {
	return &userpb.PremiumRequest{
		Id:             req.ID.String(),
		UserId:         req.UserID.String(),
		IdentityNumber: req.IdentityNumber,
		Reason:         req.Reason,
		FacePicture:    req.FacePicture,
		CreatedAt:      req.CreatedAt.Format(time.RFC3339),
	}
}

func (h *UserHandler) userToProto(user *model.User) *userpb.User {
	return &userpb.User{
		Id:             user.ID.String(),
		Name:           user.Name,
		Username:       user.Username,
		Email:          user.Email,
		Password:       user.Password,
		Gender:         user.Gender,
		ProfilePicture: user.ProfilePicture,
		Banner:         user.Banner,
		DateOfBirth:    user.DateOfBirth,
		IsActive:       user.IsActive,
		IsBanned:       user.IsBanned,
		Role:           user.Role,
		IsSubscribed:   user.IsSubscribed,
		IsPremium:      user.IsPremium,
		Bio:            user.Bio,
		CreatedAt:      user.CreatedAt.Format(time.RFC3339),
	}
}

func (h *UserHandler) userToResponse(user *model.User) *userpb.UserResponse {
	return &userpb.UserResponse{
		User: h.userToProto(user),
	}
}

// Validation functions
func validateRegisterRequest(req *userpb.RegisterRequest) error {
	// Validate name
	if err := validateName(req.Name); err != "" {
		return status.Error(codes.InvalidArgument, err)
	}

	// Validate email
	if err := validateEmail(req.Email); err != "" {
		return status.Error(codes.InvalidArgument, err)
	}

	// Validate password
	if err := validatePassword(req.Password); err != "" {
		return status.Error(codes.InvalidArgument, err)
	}

	// Validate gender
	if err := validateGender(req.Gender); err != "" {
		return status.Error(codes.InvalidArgument, err)
	}
	// Validate date of birth
	if err := validateDateOfBirth(req.DateOfBirth); err != "" {
		return status.Error(codes.InvalidArgument, err)
	}

	// Validate security question and answer
	if err := validateSecurityQuestion(req.SecurityQuestion, req.SecurityAnswer); err != "" {
		return status.Error(codes.InvalidArgument, err)
	}

	return nil
}

// Individual validation functions
func validateName(name string) string {
	if len(name) < 4 {
		return "Name must be at least 4 characters long"
	}
	
	// Only letters and spaces allowed
	nameRegex := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if !nameRegex.MatchString(name) {
		return "Name must contain only letters and spaces"
	}
	return ""
}

func validateEmail(email string) string {
	if email == "" {
		return "Email is required"
	}
	// Basic email regex validation
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if matched, _ := regexp.MatchString(emailRegex, email); !matched {
		return "Please provide a valid email address"
	}
	if len(email) > 100 {
		return "Email address is too long"
	}
	return ""
}

func validatePassword(password string) string {
	if password == "" {
		return "Password is required"
	}
	if len(password) < 8 {
		return "Password must be at least 8 characters long"
	}

	// Check for required character types
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char == '!' || char == '@' || char == '#' || char == '$' ||
			char == '%' || char == '^' || char == '&' || char == '*' ||
			char == '(' || char == ')' || char == '-' || char == '_' ||
			char == '=' || char == '+' || char == '[' || char == ']' ||
			char == '{' || char == '}' || char == '|' || char == '\\' ||
			char == ':' || char == ';' || char == '"' || char == '\'' ||
			char == '<' || char == '>' || char == ',' || char == '.' ||
			char == '?' || char == '/' || char == '~' || char == '`':
			hasSpecial = true
		}
	}

	if !hasUpper {
		return "Password must contain at least one uppercase letter"
	}
	if !hasLower {
		return "Password must contain at least one lowercase letter"
	}
	if !hasDigit {
		return "Password must contain at least one number"
	}
	if !hasSpecial {
		return "Password must contain at least one special character"
	}

	return ""
}

func validateGender(gender string) string {
	if gender == "" {
		return "Gender is required"
	}
	validGenders := []string{"male", "female"}
	for _, validGender := range validGenders {
		if strings.ToLower(gender) == validGender {
			return ""
		}
	}
	return "Gender must be one of: male, female, other, prefer-not-to-say"
}

func validateDateOfBirth(dateOfBirth string) string {
	if dateOfBirth == "" {
		return "Date of birth is required"
	}

	// Parse the date in YYYY-MM-DD format
	birthDate, err := time.Parse("2006-01-02", dateOfBirth)
	if err != nil {
		return "Date of birth must be in YYYY-MM-DD format"
	}

	// Check if the date is not in the future
	now := time.Now()
	if birthDate.After(now) {
		return "Date of birth cannot be in the future"
	}

	// Calculate age
	age := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	// Check minimum age requirement (13 years old)
	if age < 13 {
		return "You must be at least 13 years old to register"
	}

	return ""
}

func validateSecurityQuestion(question, answer string) string {
	if question == "" {
		return "Security question is required"
	}
	if answer == "" {
		return "Security answer is required"
	}
	if len(answer) < 2 {
		return "Security answer must be at least 2 characters long"
	}
	if len(answer) > 100 {
		return "Security answer is too long"
	}
	// Validate question is one of the predefined options
	validQuestions := []string{
		"What was the name of your first pet?",
		"What city were you born in?",
		"What is your favorite video game?",
		"What was the name of your first school?",
		"What was your childhood nickname?",
	}
	for _, validQuestion := range validQuestions {
		if question == validQuestion {
			return ""
		}
	}
	return "Please select a valid security question"
}
