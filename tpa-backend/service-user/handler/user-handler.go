// project/backend/service-user/handler/user_handler.go
package handler

import (
	"context"
	"errors"

	"github.com/Wenfuuu/tpa-web/service-user/model"
	// pb "github.com/Wenfuuu/tpa-web/service-user/proto"
	pb "github.com/Wenfuuu/tpa-web/pkg/proto/user"
	"github.com/Wenfuuu/tpa-web/service-user/repository"
	"github.com/Wenfuuu/tpa-web/service-user/utils"
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

func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.UserResponse, error) {
	// Validate input
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// Check if email already exists
	exists, err := h.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}
	if exists {
		return nil, status.Error(codes.AlreadyExists, "Email already registered")
	}

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
	}

	if err := h.userRepo.Create(&user); err != nil {
		return nil, status.Error(codes.Internal, "Failed to create user")
	}

	return h.userToResponse(&user), nil
}

func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Validate input
	if err := validateLoginRequest(req); err != nil {
		return nil, err
	}

	// Find user by email
	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, "Database error")
	}

	// Verify password
	if !utils.ComparePassword(user.Password, req.Password) {
		return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
	}

	// Generate JWT token
	token, err := utils.CreateJWT(user.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to generate token")
	}

	return &pb.LoginResponse{
		Token: token,
		User:  h.userToProto(user),
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
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

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
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

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
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

	return &pb.DeleteUserResponse{
		Success: true,
	}, nil
}

func (h *UserHandler) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	modelUsers, err := h.userRepo.GetAllUsers()
	if err != nil {
		return nil, status.Error(codes.Internal, "Database error")
	}

	var pbUsers []*pb.User
	for _, u := range modelUsers {
		pbUsers = append(pbUsers, h.userToProto(&u))
	}

	return &pb.GetAllUsersResponse{
		Users: pbUsers,
	}, nil
}

// Helper functions
func (h *UserHandler) userToProto(user *model.User) *pb.User {
	return &pb.User{
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
	}
}

func (h *UserHandler) userToResponse(user *model.User) *pb.UserResponse {
	return &pb.UserResponse{
		User: h.userToProto(user),
	}
}

// Validation functions
func validateRegisterRequest(req *pb.RegisterRequest) error {
	if req.Email == "" || req.Password == "" {
		return status.Error(codes.InvalidArgument, "Missing required fields")
	}
	return nil
}

func validateLoginRequest(req *pb.LoginRequest) error {
	if req.Email == "" || req.Password == "" {
		return status.Error(codes.InvalidArgument, "Email and password are required")
	}
	return nil
}