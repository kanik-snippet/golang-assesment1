package routes

import (
	"otp-auth/controllers"
	"otp-auth/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes initializes API routes
func SetupRoutes(router *gin.Engine) {
	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		api.POST("/register", controllers.RegisterUser)
		api.POST("/login", controllers.LoginUser)
		api.POST("/verify-otp", controllers.VerifyOTP)
		api.POST("/resend-otp", controllers.ResendOTP)
		api.GET("/user", middleware.AuthMiddleware(), controllers.GetUserDetails)
	}
}
