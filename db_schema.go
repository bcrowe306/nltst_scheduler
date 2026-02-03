package main

import (
	"context"
	"errors"
	"log"

	"github.com/bcrowe306/nltst_scheduler.git/models"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Helper function to check specifically for a duplicate key error
func IsDup(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	return false
}

func createCollection(database *mongo.Database, name string) *mongo.Collection {
	err := database.CreateCollection(context.TODO(), name)
	if err != nil {
		log.Fatal("Error creating collection:", err)
	}
	return database.Collection(name)
}

func createDbSchema(config *Config, database *mongo.Database) {
	usersColl := createCollection(database, "users")
	usersColl.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    map[string]interface{}{"email": 1},
		Options: options.Index().SetUnique(true),
	})

	// Create admin user
	res, err := models.CreateUser(database, "Administrator", config.AdminEmail, config.AdminPassword, "")
	if err != nil {
		if IsDup(err) {
			log.Println("Admin user already exists, skipping insertion")
		} else {
			log.Fatal("Error inserting admin user:", err)
		}
	}

	// Log admin user creation result
	if res != nil && res.InsertedID != nil {
		log.Println("Admin user created with ID:", res.InsertedID)
	}
}
