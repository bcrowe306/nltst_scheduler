package routes

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/bcrowe306/nltst_scheduler.git/models"
	"github.com/gofiber/fiber/v3/middleware/session"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Protected(c fiber.Ctx) error {
	sess := session.FromContext(c)
	if sess == nil {
		log.Print("Session is Nil")
		return c.Redirect().To("/login")
	}

	// Check if user is authenticated
	if sess.Get("authenticated") != true {
		log.Print("Not Authenticated")
		return c.Redirect().To("/login")
	}

	return c.Next()
}

func GetUserFromSession(c fiber.Ctx) (*models.User, error) {
	sess := session.FromContext(c)
	if sess == nil {
		return nil, fiber.ErrUnauthorized
	}

	userID := sess.Get("user_id")
	if userID == nil {
		return nil, fiber.ErrUnauthorized
	}

	db, ok := fiber.GetState[*mongo.Database](c.App().State(), "db")
	if !ok {
		return nil, fiber.ErrInternalServerError
	}

	user, err := models.FindUserByID(db, userID.(string))
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	return user, nil
}

func GetRoutePathList(c fiber.Ctx) []string {
	path := c.Path()

	// Split the path into segments by /
	var segments []string
	for _, segment := range strings.Split(path, "/") {
		if segment != "" {
			segments = append(segments, segment)
		}
	}
	return segments
}

func GetDefaultTemplateData(c fiber.Ctx, title string, sidebar_nav string) fiber.Map {
	user, err := GetUserFromSession(c)
	if err != nil {
		return fiber.Map{}
	}

	return fiber.Map{
		"Title":      title,
		"Path":       GetRoutePathList(c),
		"User":       user,
		"SidebarNav": sidebar_nav,
	}
}

func GetDatabaseFromContext(c fiber.Ctx) (*mongo.Database, error) {
	db, ok := fiber.GetState[*mongo.Database](c.App().State(), "db")
	if !ok {
		return nil, fiber.ErrInternalServerError
	}
	return db, nil
}

func CreateAuthRoutes(app *fiber.App) {

	app.Get("/login", func(c fiber.Ctx) error {
		return c.Render("pages/auth/login", GetDefaultTemplateData(c, "Login", ""))
	})

	app.Get("/signup", func(c fiber.Ctx) error {
		return c.Render("pages/auth/signup", GetDefaultTemplateData(c, "Sign Up", ""))
	})

	app.Post("/auth/signup", func(c fiber.Ctx) error {
		db, ok := fiber.GetState[*mongo.Database](c.App().State(), "db")
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("Database not found in context")
		}

		name := c.FormValue("name")
		email := c.FormValue("email")
		password := c.FormValue("password")
		phoneNumber := c.FormValue("phoneNumber")

		_, err := models.CreateUser(db, name, email, password, phoneNumber)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			return c.Render("pages/auth/signup", fiber.Map{
				"Title": "Sign Up",
				"Error": "Error creating user",
			})
		}

		return c.Redirect().To("/login")
	})

	app.Post("/auth/login", func(c fiber.Ctx) error {
		db, ok := fiber.GetState[*mongo.Database](c.App().State(), "db")
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("Database not found in context")
		}

		sess := session.FromContext(c)

		email := c.FormValue("email")
		password := c.FormValue("password")

		// Authenticate user
		user, err := models.FindUserByEmailPassword(db, email, password)
		if err == nil {
			// e := sess.Regenerate()
			// if e != nil {
			// 	return c.Status(500).SendString("Session error")
			// }

			if user.Enabled == false {
				return c.Render("pages/auth/login", fiber.Map{
					"Title": "Login",
					"Error": "Account disabled",
				})
			}

			models.UpdateUserLoginTime(db, user.ID)

			sess.Set("user_id", user.ID)
			sess.Set("authenticated", true)
			return c.Redirect().To("/")
		}

		return c.Render("pages/auth/login", fiber.Map{
			"Title": "Login",
			"Error": "Invalid credentials",
		})
	})

	app.Get("/auth/logout", Protected, func(c fiber.Ctx) error {
		sess := session.FromContext(c)

		// Complete session reset (clears all data + new session ID)
		if err := sess.Reset(); err != nil {
			return c.Status(500).SendString("Session error")
		}

		return c.Redirect().To("/login")
	})
}
