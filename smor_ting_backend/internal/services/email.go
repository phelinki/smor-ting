package services

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	dialer *gomail.Dialer
	from   string
}

func NewEmailService() *EmailService {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		host = "smtp.gmail.com" // Default for development
	}

	portStr := os.Getenv("SMTP_PORT")
	port := 587 // Default port
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = username
	}

	dialer := gomail.NewDialer(host, port, username, password)

	return &EmailService{
		dialer: dialer,
		from:   from,
	}
}

func (e *EmailService) SendOTP(to, otp, purpose string) error {
	subject := "Smor-Ting Verification Code"
	body := e.generateOTPEmailBody(otp, purpose)

	m := gomail.NewMessage()
	m.SetHeader("From", e.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return e.dialer.DialAndSend(m)
}

func (e *EmailService) generateOTPEmailBody(otp, purpose string) string {
	var action string
	switch purpose {
	case "registration":
		action = "complete your registration"
	case "login":
		action = "log in to your account"
	case "password_reset":
		action = "reset your password"
	default:
		action = "verify your identity"
	}

	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Smor-Ting Verification Code</title>
		<style>
			body {
				font-family: 'Poppins', Arial, sans-serif;
				line-height: 1.6;
				color: #333;
				max-width: 600px;
				margin: 0 auto;
				padding: 20px;
				background-color: #f8f9fa;
			}
			.container {
				background-color: white;
				padding: 40px;
				border-radius: 10px;
				box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
			}
			.header {
				text-align: center;
				margin-bottom: 30px;
			}
			.logo {
				color: #D21034;
				font-size: 32px;
				font-weight: bold;
				margin-bottom: 10px;
			}
			.subtitle {
				color: #002868;
				font-size: 16px;
			}
			.otp-code {
				background-color: #f8f9fa;
				border: 2px dashed #D21034;
				border-radius: 8px;
				padding: 20px;
				text-align: center;
				margin: 30px 0;
			}
			.otp-number {
				font-size: 36px;
				font-weight: bold;
				color: #D21034;
				letter-spacing: 8px;
				margin: 10px 0;
			}
			.footer {
				text-align: center;
				margin-top: 30px;
				padding-top: 20px;
				border-top: 1px solid #eee;
				color: #666;
				font-size: 14px;
			}
			.warning {
				background-color: #fff3cd;
				border: 1px solid #ffeaa7;
				border-radius: 5px;
				padding: 15px;
				margin: 20px 0;
				color: #856404;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<div class="logo">Smor-Ting</div>
				<div class="subtitle">Handyman Services for Liberia</div>
			</div>
			
			<h2>Verification Code</h2>
			<p>Hi there!</p>
			<p>Use the following verification code to %s:</p>
			
			<div class="otp-code">
				<div>Your verification code is:</div>
				<div class="otp-number">%s</div>
			</div>
			
			<div class="warning">
				<strong>Important:</strong> This code will expire in 10 minutes. Do not share this code with anyone.
			</div>
			
			<p>If you didn't request this code, please ignore this email or contact our support team.</p>
			
			<div class="footer">
				<p>Best regards,<br>The Smor-Ting Team</p>
				<p>This is an automated message, please do not reply to this email.</p>
			</div>
		</div>
	</body>
	</html>
	`, action, otp)
}