package services

import (
	"book_crud/config"
	"fmt"
	"net/smtp"
	"time"
)

type EmailService struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewEmailService() *EmailService {

	return &EmailService{
		host:     config.AppConfig.SMTPHost,
		port:     config.AppConfig.SMTPPort,
		username: config.AppConfig.SMTPUsername,
		password: config.AppConfig.SMTPPassword,
		from:     config.AppConfig.SMTPFrom,
	}

}

type EmailData struct {
	To      string
	Subject string
	Body    string
}

func (s *EmailService) SendVerificationEmail(email, username, token string) error {
	verificationURL := fmt.Sprintf("%s/api/auth/verify-email?token=%s", config.AppConfig.APPUrl, token)

	subject := "verify your email address"
	body := s.getVerificationEmailTemplate(username, verificationURL, 1)

	emaildata := EmailData{
		To:      email,
		Subject: subject,
		Body:    body,
	}

	return s.SendEmail(emaildata)
}

func (s *EmailService) SendPasswordResetEmail(email, username, token string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", config.AppConfig.APPUrl, token)

	subject := "Reset Your Password"
	body := s.getPasswordResetEmailTemplate(username, resetURL, 1)

	emailData := EmailData{
		To:      email,
		Subject: subject,
		Body:    body,
	}

	return s.SendEmail(emailData)
}

func (s *EmailService) SendEmail(data EmailData) error {
	
	//server connection string
	serverAddr := fmt.Sprintf("%s:%s", s.host, s.port)

	// Authentication
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	// construct email headers
	headers := make(map[string]string)
	headers["From"] = s.from
	headers["To"] = data.To
	headers["Subject"] = data.Subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

		// Construct message with headers
		message := ""
		for key, value := range headers {
			message += fmt.Sprintf("%s: %s\r\n", key, value)
		}
		message += "\r\n" + data.Body
	
		// Send the email
		return smtp.SendMail(
			serverAddr,
			auth,
			s.from,
			[]string{data.To},
			[]byte(message),
		)
}

func (s *EmailService) getVerificationEmailTemplate(username, verificationURL string, expiryHours int) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Verification</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            background-color: #f9f9f9;
            border-radius: 5px;
            padding: 20px;
            border: 1px solid #ddd;
        }
        .button {
            display: inline-block;
            background-color: #4CAF50;
            color: white;
            text-decoration: none;
            padding: 10px 20px;
			border-radius: 5px;
            margin: 20px 0;
        }
        .footer {
            margin-top: 20px;
            font-size: 12px;
            color: #777;
            text-align: center;
        }
    </style>
</head>
<body>
    <div class="container">
        <h2>Email Verification</h2>
        <p>Hello %s,</p>
        <p>Thank you for registering! Please verify your email address by clicking the button below:</p>
        
        <a href="%s" class="button">Verify Email Address</a>
        
        <p>This link will expire in %d hour/s.</p>
        
        <p>If the button doesn't work, copy and paste the following link into your browser:</p>
        <p>%s</p>
        
        <p>If you didn't create an account, you can safely ignore this email.</p>
    </div>
    <div class="footer">
        <p>&copy; %d SaloBook. All rights reserved.</p>
		 <p>This is an automated message, please do not reply.</p>
    </div>
</body>
</html>
`, username, verificationURL, expiryHours, verificationURL, time.Now().Year())
}

func (s *EmailService) getPasswordResetEmailTemplate(username, resetURL string, expiryHours int) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Password Reset</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            background-color: #f9f9f9;
            border-radius: 5px;
            padding: 20px;
            border: 1px solid #ddd;
        }
        .button {
            display: inline-block;
            background-color: #4CAF50;
            color: white;
            text-decoration: none;
            padding: 10px 20px;
			border-radius: 5px;
            margin: 20px 0;
        }
        .footer {
            margin-top: 20px;
            font-size: 12px;
            color: #777;
            text-align: center;
        }
    </style>
</head>
<body>
    <div class="container">
        <h2>Password Reset</h2>
        <p>Hello %s,</p>
        <p>You have requested to reset your password. Please click the button below to set a new password:</p>
        
        <a href="%s" class="button">Reset Password</a>
        
        <p>This link will expire in %d hour/s.</p>
        
        <p>If the button doesn't work, copy and paste the following link into your browser:</p>
        <p>%s</p>
        
        <p>If you didn't request a password reset, you can safely ignore this email.</p>
    </div>
    <div class="footer">
        <p>&copy; %d SaloBook. All rights reserved.</p>
		 <p>This is an automated message, please do not reply.</p>
    </div>
</body>
</html>
`, username, resetURL, expiryHours, resetURL, time.Now().Year())
}
