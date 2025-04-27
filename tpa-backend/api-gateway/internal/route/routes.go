package route

import (
	"net/http"

	"github.com/Wenfuuu/tpa-web/api-gateway/internal/handler"
	"github.com/Wenfuuu/tpa-web/api-gateway/internal/middleware"
	// pb "github.com/Wenfuuu/tpa-web/service-user/proto"
	authpb "github.com/Wenfuuu/tpa-web/pkg/proto/auth"
	userpb "github.com/Wenfuuu/tpa-web/pkg/proto/user"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

func generateCaptcha() (string, string){
	driver := base64Captcha.NewDriverDigit(100, 240, 4, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)
	_, b64s, ans, _ := captcha.Generate()
	return b64s, ans
}

func InitRoutes(
	r *gin.Engine, 
	userClient userpb.UserServiceClient,
	authClient authpb.AuthServiceClient,
) {

	userHandler := handler.NewUserHandler(userClient)
	authHandler := handler.NewAuthHandler(authClient)
	authMiddleware := middleware.NewAuthMiddleware(authClient)

	r.GET("/captcha", func(c *gin.Context) {
		captchaImage, captchaAns := generateCaptcha()
		c.JSON(http.StatusOK, gin.H{"image": captchaImage, "answer": captchaAns})
	})

	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			// users.GET("/me", authMiddleware(), getProfile)
			// users.PUT("/:id", authMiddleware(), updateUser)
			// users.DELETE("/:id", authMiddleware(), deleteUser)
			users.GET("/user", userHandler.GetAllUsers)
			users.POST("/upload/:type", userHandler.UploadImage)
		}
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/check", authMiddleware.AuthRequired())
			auth.GET("/refresh", authHandler.RefreshToken)
			auth.GET("/admin", authMiddleware.AdminRequired())
		}
	}

}