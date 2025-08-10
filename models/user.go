package models

type Users struct {
	Id         string `json:"id" example:"123"`
	UserName   string `json:"user_name" example:"newuser"`
	Email      string `json:"email" example:"newuser@example.com"`
	Password   string `json:"password" example:"yourpassword"`
	Role       string `json:"role" example:"user"`
	VerifiedAt string `json:"verified_at" example:"2023-01-01T00:00:00Z"`
}

// RegisterRequest is used for registration endpoint (Swagger and request binding)
type RegisterRequest struct {
	UserName   string `json:"user_name" example:"newuser"`
	Email      string `json:"email" example:"newuser@example.com"`
	Password   string `json:"password" example:"yourpassword"`
	Role       string `json:"role" example:"user"`
}

// LoginRequest is used for login endpoint (Swagger and request binding)
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"yourpassword"`
}

type Email struct {
	Email string `json:"email" example:"user@example.com"`
}

type TempUsers struct {
	UserName   string `json:"user_name" example:"newuser"`
	Email      string `json:"email" example:"newuser@example.com"`
	Password   string `json:"password" example:"yourpassword"`
	Role       string `json:"role" example:"user"`
	Otp        string `json:"otp" example:"123456"`
	OtpExpires string `json:"otp_expires" example:"2023-01-01T00:00:00Z"`
}

type VerifyRequest struct {
	Email string `json:"email" example:"user@example.com"`
	OTP   string `json:"otp" example:"123456"`
}

type ForgotPassword struct {
	Email      string `json:"email" example:"user@example.com"`
	OTP        string `json:"otp" example:"123456"`
	OtpExpires string `json:"otp_expires" example:"2023-01-01T00:00:00Z"`
	Verified   bool   `json:"verified"`
}

type ResetPassword struct {
	Password string `json:"password" example:"yournewpassword"`
	Email    string `json:"email" example:"user@example.com"`
}
