package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI      string
	MongoDatabase string
	AdminEmail    string
	AdminPassword string
	Port          string
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port
	}

	return &Config{
		MongoURI:      mongoURI,
		MongoDatabase: mongoDB,
		AdminEmail:    adminEmail,
		AdminPassword: adminPassword,
		Port:          port,
	}, nil
}
