package routes

import (
	"github.com/gofiber/fiber/v3"

	"log"

	"github.com/bcrowe306/nltst_scheduler.git/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func CreateUsersRoutes(app *fiber.App) {
	app.Get("/users/new", func(c fiber.Ctx) error {
		// New user page
		err := c.Render("pages/users/new", fiber.Map{
			"Title": "New User",
		}, "layouts/main")
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	app.Get("/users/:id", func(c fiber.Ctx) error {
		// User edit page
		db, ok := fiber.GetState[*mongo.Database](c.App().State(), "db")
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("Database not found in context")
		}

		userID := c.Params("id")
		user, err := models.FindUserByID(db, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching user")
		}

		err = c.Render("pages/users/edit", fiber.Map{
			"Title": "Edit User",
			"User":  user,
		}, "layouts/main")

		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	app.Get("/users", func(c fiber.Ctx) error {
		db, ok := fiber.GetState[*mongo.Database](c.App().State(), "db")
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("Database not found in context")
		}

		users, err := models.GetAllUsers(db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching users")
		}

		err = c.Render("pages/users/index", fiber.Map{
			"Title": "Users",
			"Users": users,
		}, "layouts/main")

		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	// New user handler
	app.Post("/users", func(c fiber.Ctx) error {
		db, ok := fiber.GetState[*mongo.Database](c.App().State(), "db")
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("Database not found in context")
		}

		var formData struct {
			Name        string `form:"name"`
			Email       string `form:"email"`
			PhoneNumber string `form:"phoneNumber"`
			Password    string `form:"password"`
		}

		if err := c.Bind().Form(&formData); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
		}

		if _, err := models.CreateUser(db, formData.Name, formData.Email, formData.Password, formData.PhoneNumber); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error creating user")
		} else {
			return c.Redirect().To("/users")
		}
	})

	// Update user handler
	app.Post("/users/:id", func(c fiber.Ctx) error {
		db, ok := fiber.GetState[*mongo.Database](c.App().State(), "db")
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("Database not found in context")
		}

		userID := c.Params("id")
		var updatedData models.User
		if err := c.Bind().Form(&updatedData); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
		}

		err := models.UpdateUser(db, userID, &updatedData)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error updating user")
		}

		return c.Redirect().To("/users")
	})

	// Delete user handler
	app.Get("/users/delete/:id", func(c fiber.Ctx) error {
		db, ok := fiber.GetState[*mongo.Database](c.App().State(), "db")
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("Database not found in context")
		}

		userID := c.Params("id")
		err := models.DeleteUser(db, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error deleting user")
		}

		return c.Redirect().To("/users")
	})

}
