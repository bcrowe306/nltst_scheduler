package main

import (
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3/log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/template/html/v2"

	"context"

	"github.com/bcrowe306/nltst_scheduler.git/pages"
	"github.com/bcrowe306/nltst_scheduler.git/routes"
	"github.com/bcrowe306/nltst_scheduler.git/services"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/storage/mongodb/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Render(c fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func main() {
	// Load configuration from environment variables
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}
	log.Info("Config loaded successfully")

	// Connect to MongoDB database
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

	// CONFIGURE SERVICES
	twilioService := services.NewTwilioService(config.TwilioAccountSID, config.TwilioAuthToken, config.TwilioFromNumber)
	clicksendService := services.NewClickSendService(config.ClickSendUsername, config.ClickSendAPIKey, config.ClickSendFromNumber)
	sendgridService := services.NewSendGridService(config.SendGridAPIKey)

	// Start Fiber app with HTML template engine
	engine := html.New("./views", ".html")
	engine.Reload(true) // Enable hot-reloading of templates during development

	// Create Fiber app
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Setup app state. This makes config and database accessible in handlers
	app.State().Set("config", config)
	app.State().Set("db", database)
	app.State().Set("twilioService", twilioService)
	app.State().Set("clicksendService", clicksendService)
	app.State().Set("sendgridService", sendgridService)

	// Setup session middleware with MongoDB storage
	store := mongodb.New(mongodb.Config{
		ConnectionURI: config.MongoURI,
		Database:      config.MongoDatabase,
	})
	sess := session.New(session.Config{
		Storage:         store,
		CookieSameSite:  "Lax",            // Mitigate CSRF
		IdleTimeout:     30 * time.Minute, // Session timeout
		AbsoluteTimeout: 24 * time.Hour,   // Maximum session life
		Extractor:       extractors.FromCookie("session_id"),
	})
	app.Use(sess)

	// Serve static files
	app.Use("/public", static.New("./public"))

	app.Use(routes.BreadcrumbMiddleware)

	// Create all routes
	routes.CreateAllRoutes(app)

	app.Get("/templ", func(c fiber.Ctx) {
		Render(c, pages.Index(nil))
	})

	// Start the server
	app.Listen(":" + config.Port)
}
