package routes

import (
	"otp-auth/controllers"
	"otp-auth/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/register", controllers.RegisterUser)
		api.POST("/login", controllers.LoginUser)
		api.POST("/verify-otp", controllers.VerifyOTP)
		api.POST("/resend-otp", controllers.ResendOTP)
		api.GET("/user", middleware.AuthMiddleware(), controllers.GetUserDetails)
	}
}
