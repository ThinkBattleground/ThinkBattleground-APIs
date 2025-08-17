package config

import (
	"bytes"
	"crypto/tls"
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

	from := os.Getenv("FROM_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	log.Println(from, password, smtpHost, smtpPort)

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: OTP Verification for Think Battleground\n" +
		"MIME-Version: 1.0\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\n\n" +
		emailBody

	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Connect to the SMTP server
	client, err := smtp.Dial(smtpHost + ":" + smtpPort)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	// Set up TLS configuration
	if err := client.StartTLS(&tls.Config{
		InsecureSkipVerify: true,
	}); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if err := client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to send data: %w", err)
	}
	defer wc.Close()

	if _, err := wc.Write([]byte(msg)); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}
