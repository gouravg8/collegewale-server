package email

import (
	"bytes"
	"html/template"
	"net/smtp"
	"os"
)

type EmailService struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewEmailService() *EmailService {
	return &EmailService{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}
}

func (es *EmailService) SendTemplateEmail(to, subject, templatePath string, data any) error {
	// Parse template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	// Execute template with dynamic data
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	// Prepare email
	msg := []byte("Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		body.String())

	auth := smtp.PlainAuth("", es.Username, es.Password, es.Host)
	addr := es.Host + ":" + es.Port

	return smtp.SendMail(addr, auth, es.Username, []string{to}, msg)
}
