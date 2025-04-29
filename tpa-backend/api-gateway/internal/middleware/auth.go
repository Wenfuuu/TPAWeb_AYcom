package middleware

import (
	"context"
	"net/http"

	authpb "github.com/Wenfuuu/tpa-web/pkg/proto/auth"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authClient authpb.AuthServiceClient
}

func NewAuthMiddleware(authClient authpb.AuthServiceClient) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
	}
}

func (m *AuthMiddleware) CheckGuest() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		accessTokenString := authHeader[len("Bearer "):]

		resp, err := m.authClient.CheckAuth(context.Background(), &authpb.CheckAuthRequest{
			Token: accessTokenString,
		})
		if err == nil && resp.IsAuthenticated {
			c.JSON(http.StatusForbidden, gin.H{"message": "Authenticated users are not allowed"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m* AuthMiddleware) CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header required"})
			c.Abort()
			return
		}

		accessTokenString := authHeader[len("Bearer "):]

		resp, err := m.authClient.CheckAuth(context.Background(), &authpb.CheckAuthRequest{
			Token: accessTokenString,
		})
		if err != nil || !resp.IsAuthenticated {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			c.Abort()
			return
		}

		// You can safely extract what you need
		c.Set("user_id", resp.User.Id)
		c.Set("user_role", resp.User.Role)

		c.Next()
	}
}

func (m* AuthMiddleware) CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header required"})
			c.Abort()
			return
		}

		accessTokenString := authHeader[len("Bearer "):]

		resp, err := m.authClient.AuthorizeAdmin(context.Background(), &authpb.AuthorizeAdminRequest{
			Token: accessTokenString,
		})
		if err != nil || !resp.IsAuthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token or not authorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}