package services

import (
	"encoding/json"
	"fmt"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioService struct {
	sID       string
	authToken string
	Client    *twilio.RestClient
	From      string
}

func NewTwilioService(accountSID, authToken, fromNumber string) *TwilioService {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	return &TwilioService{
		sID:       accountSID,
		authToken: authToken,
		Client:    client,
		From:      fromNumber,
	}
}

func (s *TwilioService) TSendSMS(to string, body string) error {
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(s.From)
	params.SetBody(body)

	resp, err := s.Client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("Error sending SMS message: %w", err)
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
		return nil
	}
}
