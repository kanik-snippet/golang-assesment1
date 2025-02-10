package services

import (
	"context"
	"fmt"
	"math/rand"
	"otp-auth/config"
	"time"
)

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func SaveOTP(mobile, otp string) {
	config.RedisClient.Set(context.Background(), mobile, otp, 5*time.Minute)
}

func ValidateOTP(mobile, otp string) bool {
	// Get stored OTP from Redis
	storedOTP, _ := config.RedisClient.Get(context.Background(), mobile).Result()

	// Check if OTP is valid
	if storedOTP == otp {
		// 🚀 Delete OTP from Redis after successful verification
		config.RedisClient.Del(context.Background(), mobile)
		return true
	}
	return false
}

// ✅ Check if OTP Cooldown is Active
func IsOTPCooldownActive(mobile string) (bool, time.Duration) {
	ttl, err := config.RedisClient.TTL(context.Background(), mobile+"_cooldown").Result()
	if err != nil {
		fmt.Println("Error checking OTP cooldown:", err)
		return false, 0
	}
	if ttl > 0 {
		return true, ttl
	}
	return false, 0
}

// ✅ Save OTP with Cooldown
func SaveOTPCooldown(mobile string) {
	config.RedisClient.Set(context.Background(), mobile+"_cooldown", "active", 30*time.Second) // 🔥 30 sec cooldown
}
