package models

type User struct {
	Mobile    string `json:"mobile"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Verified  bool   `json:"verified"`
	DeviceID  string `json:"device_id"`
}

type RegisterRequest struct {
	Mobile    string `json:"mobile"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type LoginRequest struct {
	Mobile string `json:"mobile"`
}

type VerifyOTPRequest struct {
	Mobile string `json:"mobile"`
	OTP    string `json:"otp"`
}

type ResendOTPRequest struct {
	Mobile string `json:"mobile"`
}
