package services

import (
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridService struct {
	APIKey string
	client *sendgrid.Client
}

func NewSendGridService(apiKey string) *SendGridService {
	if apiKey == "" {
		log.Println("SendGrid API Key is not provided. Email functionality will not work.")
	}
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridService{
		APIKey: apiKey,
		client: client,
	}
}

func (s *SendGridService) SendEmail(fromEmail, toEmail, subject, plainTextContent, htmlContent string) error {
	from := mail.NewEmail("", fromEmail)
	to := mail.NewEmail("", toEmail)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	response, err := s.client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	if response.StatusCode >= 400 {
		return fmt.Errorf("failed to send email, status code: %d, body: %s", response.StatusCode, response.Body)
	}
	return nil
}
