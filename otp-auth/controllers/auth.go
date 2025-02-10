package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"otp-auth/config"
	"otp-auth/models"
	"otp-auth/services"

	"github.com/gin-gonic/gin"
)

// ✅ Register User - Stores user details and sends OTP
func RegisterUser(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 🔥 Check if user already exists in DB
	var existingUser models.User
	if err := config.DB.Where("mobile = ?", req.Mobile).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already registered. Please login."})
		return
	}

	// 🔥 Fetch Device ID from Headers
	deviceID := c.GetHeader("X-Device-ID")
	if deviceID == "" {
		deviceID = services.GenerateDeviceID(c)
	}
	fmt.Println("Generated Device ID:", deviceID)

	// ✅ Save user in DB
	user := models.User{
		Mobile:    req.Mobile,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Verified:  false,
		DeviceID:  deviceID,
	}
	config.DB.Create(&user)

	// ✅ Send OTP (With Cooldown)
	sendOTPWithCooldown(c, req.Mobile)
}

// ✅ Login User
func LoginUser(c *gin.Context) {
	var req models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 🔥 Check if user exists in DB
	var user models.User
	if err := config.DB.Where("mobile = ?", req.Mobile).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not registered. Please register first!"})
		return
	}

	// ✅ Check if user is verified BEFORE sending OTP
	if !user.Verified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is not verified. Please verify first!"})
		return
	}

	// 🔥 If OTP not provided, send a new OTP (only if verified)
	if req.OTP == "" {
		sendOTPWithCooldown(c, req.Mobile) // ✅ Bas call karo, error check ki zaroorat nahi
		return
	}

	// 🔥 Verify OTP from Redis
	if !services.ValidateOTP(req.Mobile, req.OTP) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	// 🔥 Generate JWT Token
	token, err := services.GenerateJWT(req.Mobile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// ✅ Return User Details
	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"login_at":   services.GetCurrentTimestamp(),
	})
}

// ✅ Verify OTP
func VerifyOTP(c *gin.Context) {
	var req models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 🔥 Validate OTP from Redis
	if !services.ValidateOTP(req.Mobile, req.OTP) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	// ✅ Update user verification status in DB
	config.DB.Model(&models.User{}).Where("mobile = ?", req.Mobile).Update("verified", true)

	// ✅ Generate JWT Token
	token, _ := services.GenerateJWT(req.Mobile)
	c.JSON(http.StatusOK, gin.H{"token": token, "message": "Account verified. You can now login."})
}

// ✅ Resend OTP
func ResendOTP(c *gin.Context) {
	var req models.ResendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 🔥 Check if user exists in DB
	var user models.User
	if err := config.DB.Where("mobile = ?", req.Mobile).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not registered. Please sign up."})
		return
	}

	// ✅ Resend OTP (With Cooldown)
	sendOTPWithCooldown(c, req.Mobile)
}

// ✅ Function to send OTP with 30s cooldown (Uses Redis)
func sendOTPWithCooldown(c *gin.Context, mobile string) {
	// 🔥 Check if OTP Cooldown is Active
	isActive, remainingTime := services.IsOTPCooldownActive(mobile)
	if isActive {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": fmt.Sprintf("OTP already sent. Try again in %d seconds.", int(remainingTime.Seconds()))})
		return
	}

	// ✅ Generate & send OTP
	otp := services.GenerateOTP()
	services.SaveOTP(mobile, otp)
	services.SendOTP(mobile, otp)

	// ✅ Save OTP Cooldown in Redis
	services.SaveOTPCooldown(mobile)

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// ✅ Get User Details
func GetUserDetails(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Extract Token
	tokenString := strings.Split(authHeader, " ")
	if len(tokenString) < 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Parse Token
	claims, err := services.ParseJWT(tokenString[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Get Mobile from Token
	mobile, ok := claims["mobile"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Fetch User Details from DB
	var user models.User
	if err := config.DB.Where("mobile = ?", mobile).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return User Details
	c.JSON(http.StatusOK, gin.H{"user": user})
}
