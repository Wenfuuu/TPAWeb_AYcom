package route

import (
	"net/http"

	"github.com/Wenfuuu/tpa-web/api-gateway/internal/handler"
	pb "github.com/Wenfuuu/tpa-web/api-gateway/proto/user"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

func generateCaptcha() (string, string){
	driver := base64Captcha.NewDriverDigit(100, 240, 4, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)
	_, b64s, ans, _ := captcha.Generate()
	return b64s, ans
}

func InitRoutes(r *gin.Engine, userClient pb.UserServiceClient) {

	userHandler := handler.NewUserHandler(userClient)

	r.GET("/captcha", func(c *gin.Context) {
		captchaImage, captchaAns := generateCaptcha()
		c.JSON(http.StatusOK, gin.H{"image": captchaImage, "answer": captchaAns})
	})

	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/register", userHandler.RegisterUser)
			users.POST("/login", userHandler.LoginUser)
			// users.GET("/me", authMiddleware(), getProfile)
			// users.PUT("/:id", authMiddleware(), updateUser)
			// users.DELETE("/:id", authMiddleware(), deleteUser)
			users.GET("/user", userHandler.GetAllUsers)
		}
	}

}