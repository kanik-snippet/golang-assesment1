package services

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

// GenerateDeviceID - Automatically generates a unique Device ID
func GenerateDeviceID(c *gin.Context) string {
	userAgent := c.GetHeader("User-Agent")
	ip := c.ClientIP()

	// 🔥 Create a unique fingerprint
	hash := sha256.New()
	hash.Write([]byte(userAgent + ip))
	return hex.EncodeToString(hash.Sum(nil))[:16] // Take first 16 chars
}
