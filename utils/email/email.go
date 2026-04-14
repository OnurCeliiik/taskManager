package email

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailService struct {
	username  string
	fromEmail string
	password  string
	smtpHost  string
	smtpPort  string
}

func NewEmailService() *EmailService {
	return &EmailService{
		fromEmail: os.Getenv("SMTP_FROM_EMAIL"),
		password:  os.Getenv("SMTP_PASSWORD"),
		smtpHost:  os.Getenv("SMTP_HOST"),
		smtpPort:  os.Getenv("SMTP_PORT"),
		username:  os.Getenv("SMTP_USERNAME"),
	}
}

func (es *EmailService) SendPasswordResetEmail(toEmail, resetToken string) error {
	// In development, skip sending if not configured
	if es.smtpHost == "" {
		return nil
	}

	subject := "Password Reset Request"

	// Build reset link (replace with your domain)
	resetLink := fmt.Sprintf("https://yourdomain.com/reset-password?token=%s", resetToken)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; }
        .container { max-width: 600px; margin: 0 auto; }
        .button {
            display: inline-block;
            padding: 10px 20px;
            background-color: #007bff;
            color: white;
            text-decoration: none;
            border-radius: 5px;
        }
        .footer { margin-top: 30px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Password Reset Request</h2>
        <p>We received a request to reset your password. Click the link below to proceed:</p>
        <p>
            <a href="%s" class="button">Reset Password</a>
        </p>
        <p>Or copy and paste this link in your browser:</p>
        <p>%s</p>
        <p>This link expires in 1 hour.</p>
        <p>If you didn't request this, you can safely ignore this email.</p>
        <div class="footer">
            <p>Task Manager API</p>
        </div>
    </div>
</body>
</html>
    `, resetLink, resetLink)

	return es.sendEmail(toEmail, subject, body)
}

func (es *EmailService) SendRegistrationEmail(toEmail, userName string) error {
	if es.smtpHost == "" {
		return nil
	}

	subject := "Welcome to Task Manager"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; }
        .container { max-width: 600px; margin: 0 auto; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Welcome, %s!</h2>
        <p>Your account has been successfully created.</p>
        <p>You can now log in to Task Manager and start managing your tasks.</p>
        <p>Happy tasks managing!</p>
    </div>
</body>
</html>
    `, userName)

	return es.sendEmail(toEmail, subject, body)
}

func (es *EmailService) sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", es.username, es.password, es.smtpHost)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		es.fromEmail, to, subject, body)

	addr := fmt.Sprintf("%s:%s", es.smtpHost, es.smtpPort)
	err := smtp.SendMail(addr, auth, es.fromEmail, []string{to}, []byte(msg))

	return err
}
