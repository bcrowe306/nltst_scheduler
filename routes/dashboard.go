package routes

import (
	"github.com/gofiber/fiber/v3"

	"log"
)

func CreateDashboardRoutes(app *fiber.App) {
	app.Get("/", Protected, func(c fiber.Ctx) error {

		err := c.Render("pages/dashboard/index", GetDefaultTemplateData(c, "Dashboard"), "layouts/main")

		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}

		return nil
	})
}
