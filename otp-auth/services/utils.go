package services

import (
	"time"
)

// ✅ GetCurrentTimestamp - Returns current timestamp in UNIX format
func GetCurrentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
