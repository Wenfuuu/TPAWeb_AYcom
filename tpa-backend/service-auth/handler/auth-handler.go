package handler

import (
	"context"

	authpb "github.com/Wenfuuu/tpa-web/pkg/proto/auth"
	userpb "github.com/Wenfuuu/tpa-web/pkg/proto/user"
	"github.com/Wenfuuu/tpa-web/service-auth/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	userClient userpb.UserServiceClient
}

func NewAuthHandler(userClient userpb.UserServiceClient) *AuthHandler {
	return &AuthHandler{
		userClient: userClient,
	}
}

func (h *AuthHandler) Login(ctx context.Context, req *authpb.AuthLoginRequest) (*authpb.AuthLoginResponse, error) {
	// Find user by email
	userResp, err := h.userClient.GetUserByEmail(ctx, &userpb.GetUserByEmailRequest{Email: req.Email})
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	// Verify password
	if !utils.ComparePassword(userResp.User.Password, req.Password) {
		return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
	}

	// Create token
	accessToken, refreshToken, err := utils.CreateJWT(userResp.User.Id, userResp.User.Role)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create token")
	}

	return &authpb.AuthLoginResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		User: h.userToProto(userResp.User),
	}, nil
}

func (h *AuthHandler) Register(ctx context.Context, req *authpb.AuthRegisterRequest) (*authpb.AuthRegisterResponse, error) {
	// Hash password
	hashedPassword, err := utils.EncryptPassword(req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to hash password")
	}

	// Call User Service to create user
	createResp, err := h.userClient.Register(ctx, &userpb.RegisterRequest{
		Name:           req.Name,
		Username:       req.Username,
		Email:          req.Email,
		Password:       hashedPassword,
		Gender:         req.Gender,
		ProfilePicture: req.ProfilePicture,
		Banner:         req.Banner,
		DateOfBirth:    req.DateOfBirth,
		IsActive:       req.IsActive,
		IsBanned:       req.IsBanned,
		Role:           req.Role,
		IsSubscribed: 	req.IsSubscribed,
		IsPremium: 		req.IsPremium,
		Bio: 			req.Bio,
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Failed to create user")
	}

	return &authpb.AuthRegisterResponse{
		User: h.userToProto(createResp.User),
	}, nil
}

func (h *AuthHandler) CheckAuth(ctx context.Context, req *authpb.CheckAuthRequest) (*authpb.CheckAuthResponse, error) {
	user_id, _, err := utils.ParseAccessToken(req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid access token")
	}

	user, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{
		Id: user_id,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get user info")
	}

	return &authpb.CheckAuthResponse{
		IsAuthenticated: true,
		User: &authpb.AuthUser{
			Id:             user.User.Id,
			Name:           user.User.Name,
			Username:       user.User.Username,
			Email:          user.User.Email,
			Password:       user.User.Password,
			Gender:         user.User.Gender,
			ProfilePicture: user.User.ProfilePicture,
			Banner:         user.User.Banner,
			DateOfBirth:    user.User.DateOfBirth,
			IsActive:       user.User.IsActive,
			IsBanned:       user.User.IsBanned,
			Role:           user.User.Role,
			IsSubscribed:   user.User.IsSubscribed,
			IsPremium:      user.User.IsPremium,
			Bio:            user.User.Bio,
			CreatedAt:      user.User.CreatedAt,
		},
	}, nil
}

func (h *AuthHandler) AuthorizeAdmin(ctx context.Context, req *authpb.AuthorizeAdminRequest) (*authpb.AuthorizeAdminResponse, error) {
	user_id, _, err := utils.ParseAccessToken(req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid access token")
	}

	user, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{
		Id: user_id,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get user info")
	}

	if user.User.Role == "user" {
		return nil, status.Error(codes.Unauthenticated, "You are not authorized to access this page")
	}

	return &authpb.AuthorizeAdminResponse{
		IsAuthorized: true,
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.AuthLoginResponse, error) {
	// Parse refresh token
	userID, role, err := utils.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid refresh token")
	}

	// Get user info again
	userResp, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get user info")
	}

	// Create new tokens
	accessToken, newRefreshToken, err := utils.CreateJWT(userID, role)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create new tokens")
	}

	return &authpb.AuthLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         h.userToProto(userResp.User),
	}, nil
}

// Helper functions
func (h *AuthHandler) userToProto(user *userpb.User) *authpb.AuthUser {
	return &authpb.AuthUser{
		Id:        		user.Id,
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
		CreatedAt: 		user.CreatedAt,
	}
}