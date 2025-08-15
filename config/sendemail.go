package config

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"text/template"
	"thinkbattleground-apis/constants"
)

type EmailData struct {
	Username string
	OTP      string
}

// Send OTP email
func SendEmail(to, otp, username, message string) error {
	if err := LoadEnv(); err != nil {
		log.Println(constants.LOAD_ENV_ERROR)
		return err
	}

	templates, err := template.ParseFiles("templates/send_registration_otp.html")
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	data := struct {
		Username string
		OTP      string
		Message  string
	}{
		Username: username,
		OTP:      otp,
		Message:  message,
	}

	// Execute template
	var body bytes.Buffer
	if err := templates.Execute(&body, data); err != nil {
		log.Println(err)
		panic(err)
	}

	emailBody := body.String()

	from := os.Getenv("EMAIL_USER_FROM")
	password := os.Getenv("EMAIL_USER_PASSWORD")
	smtpHost := os.Getenv("EMAIL_SMTP_HOST")
	smtpPort := os.Getenv("EMAIL_SMTP_PORT")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: OTP Verification for Think Battleground\n" +
		"MIME-Version: 1.0\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\n\n" +
		emailBody

	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
}
