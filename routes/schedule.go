package routes

import (
	"github.com/gofiber/fiber/v3"

	"log"
)

func CreateScheduleRoutes(app *fiber.App) {
	app.Get("/schedule", Protected, func(c fiber.Ctx) error {

		err := c.Render("pages/schedule/index", GetDefaultTemplateData(c, "Schedule", "schedule"), "layouts/main")

		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})
}
