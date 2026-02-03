package routes

import (
	"github.com/gofiber/fiber/v3"

	"log"

	"github.com/bcrowe306/nltst_scheduler.git/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func CreateMembersRoutes(app *fiber.App) {
	app.Get("/members", func(c fiber.Ctx) error {
		db, ok := fiber.GetState[*mongo.Database](c.App().State(), "db")
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("Database not found in context")
		}

		users, err := models.GetAllUsers(db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching users")
		}

		err = c.Render("pages/members/index", fiber.Map{
			"Title": "Members",
			"Users": users,
		}, "layouts/main")

		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})
}
