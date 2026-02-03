package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/template/html/v2"

	"context"

	"github.com/bcrowe306/nltst_scheduler.git/routes"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// TODO: Create a config struct that holds all configuration values and load it from environment variables or a config file.

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	client, err := mongo.Connect(options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	defer client.Disconnect(context.TODO())

	database := client.Database(config.MongoDatabase)

	// Test the connection
	err = client.Ping(nil, nil)
	if err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}

	// Create DB schema
	createDbSchema(config, database)

	// Start Fiber app with HTML template engine
	engine := html.New("./views", ".html")
	engine.Reload(true) // Enable hot-reloading of templates during development

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.State().Set("config", config)
	app.State().Set("db", database)

	// Serve static files
	app.Use("/public", static.New("./public"))

	// Create all routes
	routes.CreateAllRoutes(app)

	// Start the server
	app.Listen(":" + config.Port)
}
