package handler

import (
	"context"
	"encoding/json"

	// "fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Wenfuuu/tpa-web/pkg/messaging"
	authpb "github.com/Wenfuuu/tpa-web/pkg/proto/auth"
	userpb "github.com/Wenfuuu/tpa-web/pkg/proto/user"
	"github.com/Wenfuuu/tpa-web/service-auth/utils"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	// "gopkg.in/gomail.v2"
)

type AuthHandler struct {
	userClient    userpb.UserServiceClient
	rdb           *redis.Client
	rabbitChannel *amqp091.Channel
}

func NewAuthHandler(userClient userpb.UserServiceClient) *AuthHandler {
	ch, err := utils.GetRabbitMQChannel()
	if err != nil {
		log.Fatalf("Failed to get RabbitMQ channel for AuthHandler: %v", err)
	}

	exchangeName := "email_exchange"
	err = ch.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	return &AuthHandler{
		userClient:    userClient,
		rdb:           createRedisClient(),
		rabbitChannel: ch,
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

	// check activated
	if !userResp.GetUser().IsActive {
		return nil, status.Error(codes.Unauthenticated, "Account not activated")
	}

	// check banned
	if userResp.GetUser().IsBanned {
		return nil, status.Error(codes.Unauthenticated, "Account is banned")
	}

	// Create token
	accessToken, refreshToken, err := utils.CreateJWT(userResp.User.Id, userResp.User.Role)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create token")
	}

	return &authpb.AuthLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         h.userToProto(userResp.User),
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
		Name:             req.Name,
		Username:         req.Username,
		Email:            req.Email,
		Password:         hashedPassword,
		Gender:           req.Gender,
		ProfilePicture:   req.ProfilePicture,
		Banner:           req.Banner,
		DateOfBirth:      req.DateOfBirth,
		IsActive:         req.IsActive,
		IsBanned:         req.IsBanned,
		Role:             req.Role,
		SecurityQuestion: req.SecurityQuestion,
		SecurityAnswer:   req.SecurityAnswer,
		IsSubscribed:     req.IsSubscribed,
		IsPremium:        req.IsPremium,
		Bio:              req.Bio,
	})
	if err != nil {
		// Check if this is a validation error or other gRPC error
		if st, ok := status.FromError(err); ok {
			// Propagate the specific error code and message from user service
			return nil, status.Error(st.Code(), st.Message())
		}
		// Fallback for non-gRPC errors
		return nil, status.Error(codes.Internal, "Failed to create user")
	}

	// Rabbitmq to email service to send success message
	payload := messaging.EmailPayload{
		Email: req.Email,
		Word:  "Successfully registered account for AY.com",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal email payload: %v", err)
		return nil, status.Error(codes.Internal, "Failed to process request (marshal)")
	}

	exchangeName := "email_exchange"
	routingKey := "email"
	err = h.rabbitChannel.PublishWithContext(
		context.Background(),
		exchangeName,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp091.Transient,
		},
	)
	if err != nil {
		log.Printf("Failed to publish email message: %v", err)
		return nil, status.Error(codes.Internal, "Failed to queue email request")
	}
	log.Printf("Successfully published email request for %s", req.Email)

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

func (h *AuthHandler) GetAuthUser(ctx context.Context, req *authpb.GetAuthUserRequest) (*authpb.AuthLoginResponse, error) {
	// parse access
	userID, _, err := utils.ParseAccessToken(req.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid access token")
	}

	userResp, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get user info")
	}

	return &authpb.AuthLoginResponse{
		User: h.userToProto(userResp.User),
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

func (h *AuthHandler) SendVerificationCode(ctx context.Context, req *authpb.SendVerificationCodeRequest) (*authpb.SendVerificationCodeResponse, error) {
	code := generateRandomCode(5)

	err := h.rdb.Set(ctx, "verify:"+req.Email, code, 5*time.Minute).Err()
	if err != nil {
		return &authpb.SendVerificationCodeResponse{
			Success: false,
		}, err
	}

	// rabbitmq logic here
	payload := messaging.VerificationEmailPayload{
		Email: req.Email,
		Code:  code,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal verification email payload: %v", err)
		return nil, status.Error(codes.Internal, "Failed to process request (marshal)")
	}

	exchangeName := "email_exchange"
	routingKey := "verification_email"
	err = h.rabbitChannel.PublishWithContext(
		context.Background(),
		exchangeName,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp091.Persistent, // Make message persistent (survive broker restart)
		},
	)
	if err != nil {
		log.Printf("Failed to publish verification email message: %v", err)
		return nil, status.Error(codes.Internal, "Failed to queue email request")
	}
	log.Printf("Successfully published verification email request for %s", req.Email)

	return &authpb.SendVerificationCodeResponse{
		Success: true,
	}, nil
}

func (h *AuthHandler) VerifyCode(ctx context.Context, req *authpb.VerifyCodeRequest) (*authpb.VerifyCodeResponse, error) {
	storedCode, err := h.rdb.Get(ctx, "verify:"+req.Email).Result()
	if err == redis.Nil {
		return &authpb.VerifyCodeResponse{
			Verified: false,
		}, nil
	} else if err != nil {
		return nil, err
	}

	if storedCode == req.Code {
		return &authpb.VerifyCodeResponse{
			Verified: true,
		}, nil
	} else {
		return nil, status.Error(codes.Unauthenticated, "Invalid verification code")
	}
}

// Helper functions
func (h *AuthHandler) userToProto(user *userpb.User) *authpb.AuthUser {
	return &authpb.AuthUser{
		Id:             user.Id,
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
		CreatedAt:      user.CreatedAt,
	}
}

func createRedisClient() *redis.Client {
	redisAddr := utils.GetDotEnv("REDIS_ADDR")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	return rdb
}

func generateRandomCode(length int) string {
	rand.Seed(time.Now().UnixNano()) // important to avoid repeating numbers
	digits := []rune("0123456789")
	code := make([]rune, length)
	for i := range code {
		code[i] = digits[rand.Intn(len(digits))]
	}
	return string(code)
}
