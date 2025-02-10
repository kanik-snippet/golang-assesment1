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

// RegisterUser godoc
// @Summary      Register a new user
// @Description  Stores user details and sends an OTP for verification
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body models.RegisterRequest true "User Registration Data"
// @Success      200  {object}  models.RegisterResponse  "OTP sent successfully"
// @Failure      400  {object}  map[string]string  "error: Invalid input"
// @Failure      409  {object}  map[string]string  "error: User already registered"
// @Failure      429  {object}  map[string]string  "error: OTP already sent. Try again later"
// @Failure      500  {object}  map[string]string  "error: Failed to register user"
// @Router       /api/register [post]
func RegisterUser(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var existingUser models.User
	if err := config.DB.Where("mobile = ?", req.Mobile).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already registered. Please login."})
		return
	}

	deviceID := c.GetHeader("X-Device-ID")
	if deviceID == "" {
		deviceID = services.GenerateDeviceID(c)
	}
	fmt.Println("Generated Device ID:", deviceID)

	user := models.User{
		Mobile:    req.Mobile,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Verified:  false,
		DeviceID:  deviceID,
	}
	config.DB.Create(&user)

	sendOTPWithCooldown(c, req.Mobile)
}

// LoginUser godoc
// @Summary      Login user
// @Description  Logs in a user and sends an OTP for verification
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body models.VerifyOTPRequest true "User Login Data"
// @Success      200  {object}  map[string]string  "token: JWT Token"
// @Failure      400  {object}  map[string]string  "error: Invalid input"
// @Failure      404  {object}  map[string]string  "error: User not found"
// @Failure      401  {object}  map[string]string  "error: Unauthorized"
// @Router       /api/login [post]
func LoginUser(c *gin.Context) {
	var req models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := config.DB.Where("mobile = ?", req.Mobile).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not registered. Please register first!"})
		return
	}

	if !user.Verified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is not verified. Please verify first!"})
		return
	}

	if req.OTP == "" {
		sendOTPWithCooldown(c, req.Mobile)
		return
	}

	if !services.ValidateOTP(req.Mobile, req.OTP) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	token, err := services.GenerateJWT(req.Mobile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"login_at":   services.GetCurrentTimestamp(),
	})
}

// VerifyOTP godoc
// @Summary      Verify OTP
// @Description  Verifies the OTP sent to the user
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body models.VerifyOTPRequest true "OTP Verification Data"
// @Success      200  {object}  map[string]string  "message: Account verified"
// @Failure      400  {object}  map[string]string  "error: Invalid input"
// @Failure      401  {object}  map[string]string  "error: Invalid or expired OTP"
// @Router       /api/verify-otp [post]
func VerifyOTP(c *gin.Context) {
	var req models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if !services.ValidateOTP(req.Mobile, req.OTP) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	config.DB.Model(&models.User{}).Where("mobile = ?", req.Mobile).Update("verified", true)

	token, _ := services.GenerateJWT(req.Mobile)
	c.JSON(http.StatusOK, gin.H{"token": token, "message": "Account verified. You can now login."})
}

// ResendOTP godoc
// @Summary      Resend OTP
// @Description  Resends OTP to the registered user
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body models.ResendOTPRequest true "Mobile Number"
// @Success      200  {object}  map[string]string  "message: OTP sent successfully"
// @Failure      400  {object}  map[string]string  "error: Invalid input"
// @Failure      404  {object}  map[string]string  "error: User not found"
// @Router       /api/resend-otp [post]
func ResendOTP(c *gin.Context) {
	var req models.ResendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := config.DB.Where("mobile = ?", req.Mobile).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not registered. Please sign up."})
		return
	}

	sendOTPWithCooldown(c, req.Mobile)
}

// GetUserDetails godoc
// @Summary      Get user details
// @Description  Fetches user details after authentication
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Success      200  {object}  models.User "User Details"
// @Failure      401  {object}  map[string]string  "error: Unauthorized"
// @Failure      404  {object}  map[string]string  "error: User not found"
// @Router       /api/user [get]
func GetUserDetails(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tokenString := strings.Split(authHeader, " ")
	if len(tokenString) < 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, err := services.ParseJWT(tokenString[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	mobile, ok := claims["mobile"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var user models.User
	if err := config.DB.Where("mobile = ?", mobile).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
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
