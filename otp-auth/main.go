package main

import (
	"otp-auth/config"
	_ "otp-auth/docs" // Import the generated Swagger docs
	"otp-auth/routes"

	"github.com/gin-gonic/gin"
)

// @title OTP Authentication API
// @version 1.0
// @description This is an OTP authentication API built with Go and Gin.
// @host localhost:8080
// @BasePath /api
func main() {
	config.ConnectDB()
	config.ConnectRedis()

	router := gin.Default()
	routes.SetupRoutes(router)

	router.Run(":8080") // Start server
}
