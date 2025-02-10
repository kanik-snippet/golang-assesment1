package main

import (
	"otp-auth/config"
	"otp-auth/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()
	config.ConnectRedis()

	router := gin.Default()
	routes.SetupRoutes(router)
	router.Run(":8080")
}
