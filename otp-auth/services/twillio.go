package services

import (
	"fmt"
	"os"
	"strings"

	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendOTP(mobile, otp string) error {
	client := twilio.NewRestClient()
	// Ensure the number is in E.164 format
	if !strings.HasPrefix(mobile, "+") {
		mobile = "+91" + mobile
	}
	// Print OTP in Terminal
	fmt.Println("Generated OTP:", otp)

	params := &openapi.CreateMessageParams{}
	params.SetTo(mobile)
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetBody("Your OTP is: " + otp)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending OTP:", err)
		return err
	}

	fmt.Println("OTP sent successfully to", mobile)
	return nil
}
