package config

import (
	"crypto/rand"
	"log"
)

// Generate random OTP
func GenerateOTP() string {
	const otpLength = 6
	const charset = "0123456789"

	b := make([]byte, otpLength)

	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Failed to generate OTP: %v", err)
	}
	
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}
