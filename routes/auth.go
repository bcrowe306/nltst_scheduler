package routes

import (
	"log"

	"github.com/bcrowe306/nltst_scheduler.git/models"
	"github.com/bcrowe306/nltst_scheduler.git/pages"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// CreateAuthRoutes sets up all authentication related routes (login, signup, logout)
func CreateAuthRoutes(app *fiber.App, BaseRoute string) {

	app.Get("/login", func(c fiber.Ctx) error {
		return RenderFullPage(c, pages.LoginPage())
		// return c.Render("pages/auth/login", GetDefaultTemplateData(c, "Login", ""))
	})

	app.Get("/signup", func(c fiber.Ctx) error {
		return RenderFullPage(c, pages.SignupPage())
	})

	app.Post(BaseRoute+"/signup", func(c fiber.Ctx) error {
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

	app.Post(BaseRoute+"/login", func(c fiber.Ctx) error {
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

	app.Get(BaseRoute+"/logout", Protected, func(c fiber.Ctx) error {
		sess := session.FromContext(c)

		// Complete session reset (clears all data + new session ID)
		if err := sess.Reset(); err != nil {
			return c.Status(500).SendString("Session error")
		}

		return c.Redirect().To("/login")
	})
}
