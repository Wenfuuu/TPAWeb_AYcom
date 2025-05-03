// project/backend/service-user/handler/user_handler.go
package handler

import (
	"bytes"
	"context"
	"errors"
	"path/filepath"
	"strings"
	"time"

	"github.com/Wenfuuu/tpa-web/service-user/model"
	"github.com/Wenfuuu/tpa-web/service-user/utils"

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
		FollowerId: follow.FollowerID.String(),
		FollowingId: follow.FollowingID.String(),
	}
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
		IsActive: 		req.IsActive,
		IsBanned: 		req.IsBanned,
		Role: 			req.Role,
		IsSubscribed: 	req.IsSubscribed,
		IsPremium: 		req.IsPremium,
		Bio: 			req.Bio,
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

	fileExt := strings.ToLower(filepath.Ext(req.Filename))
    
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
	_, err = supabase.Storage.UploadFile(bucketName, req.Filename, fileReader, storage_go.FileOptions{
		CacheControl: stringPtr("3600"),
		ContentType: &contentType,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to upload image")
	}
	
	// Get public URL
	publicURL := supabase.Storage.GetPublicUrl(bucketName, req.Filename).SignedURL
	
	return &userpb.UploadImageResponse{
		ImageUrl: publicURL,
	}, nil
}

func stringPtr(s string) *string {
    return &s
}

// Helper functions
func (h *UserHandler) userToProto(user *model.User) *userpb.User {
	return &userpb.User{
		Id:        		user.ID.String(),
		Name: 			user.Name,
		Username: 		user.Username,
		Email:     		user.Email,
		Password:     	user.Password,
		Gender: 		user.Gender,
		ProfilePicture: user.ProfilePicture,
		Banner: 		user.Banner,
		DateOfBirth: 	user.DateOfBirth,
		IsActive: 		user.IsActive,
		IsBanned: 		user.IsBanned,
		Role: 			user.Role,
		IsSubscribed: 	user.IsSubscribed,
		IsPremium: 		user.IsPremium,
		Bio: 			user.Bio,
		CreatedAt: 		user.CreatedAt.Format(time.RFC3339),
	}
}

func (h *UserHandler) userToResponse(user *model.User) *userpb.UserResponse {
	return &userpb.UserResponse{
		User: h.userToProto(user),
	}
}

// Validation functions
func validateRegisterRequest(req *userpb.RegisterRequest) error {
	if req.Email == "" || req.Password == "" {
		return status.Error(codes.InvalidArgument, "Missing required fields")
	}
	return nil
}