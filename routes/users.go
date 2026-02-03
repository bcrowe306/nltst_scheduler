package routes

import (
	"github.com/gofiber/fiber/v3"

	"log"

	"github.com/bcrowe306/nltst_scheduler.git/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func CreateUsersRoutes(app *fiber.App) {
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

}
