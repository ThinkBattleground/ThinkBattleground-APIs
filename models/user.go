package models

type Users struct {
	Id         string `json:"id"`
	UserId     string `json:"user_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	VerifiedAt string `json:"verified_at"`
}

type Email struct {
	Email string `json:"email"`
}

type TempUsers struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	Otp        string `json:"otp"`
	OtpExpires string `json:"otp_expires"`
}

type VerifyRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type ForgotPassword struct {
	Email      string `json:"email"`
	OTP        string `json:"otp"`
	OtpExpires string `json:"otp_expires"`
	Verified   bool   `json:"verified"`
}

type ResetPassword struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
