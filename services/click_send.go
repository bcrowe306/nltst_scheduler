package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CSMessage struct {
	Source string `json:"source"`
	Body   string `json:"body"`
	To     string `json:"to"`
}

type ClickSendService struct {
	Username   string
	APIKey     string
	BaseURL    string
	FromNumber string
	client     *http.Client
}

// NewClickSendService creates a new instance of ClickSendService
func NewClickSendService(username, apiKey, fromNumber string) *ClickSendService {
	return &ClickSendService{
		Username:   username,
		APIKey:     apiKey,
		BaseURL:    "https://rest.clicksend.com/v3",
		FromNumber: fromNumber,
		client:     &http.Client{},
	}
}

func (s *ClickSendService) SendSMSCS(to string, body string) error {
	reqUrl := s.BaseURL + "/sms/send"
	messages := []CSMessage{
		{
			Source: s.FromNumber,
			Body:   body,
			To:     to,
		},
	}

	payload, err := json.Marshal(map[string]interface{}{
		"messages": messages,
	})
	if err != nil {
		return fmt.Errorf("Error marshaling JSON payload: %w", err)
	}

	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("Error creating HTTP request: %w", err)
	}

	req.SetBasicAuth(s.Username, s.APIKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending HTTP request: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Error reading response body: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("Error response from ClickSend: %s", string(bodyBytes))
	}

	fmt.Println("Response from ClickSend:", string(bodyBytes))
	return nil
}
