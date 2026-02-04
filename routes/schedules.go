package routes

import (
	"github.com/gofiber/fiber/v3"

	"log"
)

func CreateSchedulesRoutes(app *fiber.App) {
	app.Get("/schedules", Protected, func(c fiber.Ctx) error {

		err := c.Render("pages/schedules/index", GetDefaultTemplateData(c, "Schedules"), "layouts/main")

		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})
}
