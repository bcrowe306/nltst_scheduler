package routes

import (
	"github.com/bcrowe306/nltst_scheduler.git/models"
	"github.com/gofiber/fiber/v3"

	"log"
)

func CreateDashboardRoutes(app *fiber.App) {
	app.Get("/", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		data := GetDefaultTemplateData(c, "Dashboard", "dashboard")
		eventsByPosition, err := models.GetEventsGroupedByPositionName(db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching events")
		}
		data["Positions"] = eventsByPosition
		err = c.Render("pages/dashboard/index", data, "layouts/main")

		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}

		return nil
	})
}
