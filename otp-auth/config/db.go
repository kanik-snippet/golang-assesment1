package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Get database credentials from .env
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// DSN (Data Source Name) string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Open GORM DB connection
	var dbErr error
	DB, dbErr = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if dbErr != nil {
		fmt.Println("❌ Database connection failed:", dbErr)
		return
	}

	fmt.Println("✅ Database connected successfully!")
}
