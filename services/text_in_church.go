package services

import "net/http"

type TextInChurchService struct {
	BaseURL string
	APIKey  string
	client  *http.Client
}
