package middleware

// import (
// 	"net/http"

// 	"github.com/Wenfuuu/tpa-web/api-gateway/utils"
// 	"github.com/gin-gonic/gin"
// )

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header required"})
// 			c.Abort()
// 			return
// 		}

// 		tokenString := authHeader[len("Bearer "):]

// 		user_id, err := utils.ParseJWT(tokenString)
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token", "error": err.Error()})
// 			c.Abort()
// 			return
// 		}

// 		// You can safely extract what you need
// 		c.Set("user_id", user_id)

// 		c.Next()
// 	}
// }