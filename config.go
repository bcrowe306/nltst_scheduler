package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI            string
	MongoDatabase       string
	AdminEmail          string
	AdminPassword       string
	TwilioAccountSID    string
	TwilioAuthToken     string
	TwilioFromNumber    string
	ClickSendAPIKey     string
	ClickSendUsername   string
	ClickSendPassword   string
	ClickSendFromNumber string
	Port                string
}

func LoadConfig() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, proceeding with system environment variables")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGODB_URI not set in environment")
	}

	mongoDB := os.Getenv("MONGODB_DATABASE")
	if mongoDB == "" {
		return nil, fmt.Errorf("MONGODB_DATABASE not set in environment")
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminEmail == "" || adminPassword == "" {
		return nil, fmt.Errorf("ADMIN_EMAIL or ADMIN_PASSWORD not set in environment")
	}

	// Twilio Credentials
	twilioAccountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	twilioFromNumber := os.Getenv("TWILIO_FROM_NUMBER")
	if twilioAccountSID == "" || twilioAuthToken == "" || twilioFromNumber == "" {
		log.Print("Twilio Credentials are not supplied. SMS will not work")
	}

	// ClickSend Credentials
	clicksendAPIKey := os.Getenv("CLICK_SEND_API_KEY")
	clicksendUsername := os.Getenv("CLICK_SEND_USERNAME")
	clicksendFromNumber := os.Getenv("CLICK_SEND_FROM_NUMBER")
	if clicksendAPIKey == "" || clicksendUsername == "" || clicksendFromNumber == "" {
		log.Print("ClickSend Credentials are not supplied. SMS will not work")
	}

	// Application Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port
	}

	return &Config{
		MongoURI:            mongoURI,
		MongoDatabase:       mongoDB,
		AdminEmail:          adminEmail,
		AdminPassword:       adminPassword,
		TwilioAccountSID:    twilioAccountSID,
		TwilioAuthToken:     twilioAuthToken,
		TwilioFromNumber:    twilioFromNumber,
		ClickSendAPIKey:     clicksendAPIKey,
		ClickSendUsername:   clicksendUsername,
		ClickSendFromNumber: clicksendFromNumber,
		Port:                port,
	}, nil
}
