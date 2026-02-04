package routes

import (
	"github.com/gofiber/fiber/v3"

	"log"
)

func CreateSettingsRoutes(app *fiber.App) {
	app.Get("/settings", Protected, func(c fiber.Ctx) error {

		err := c.Render("pages/settings/index", GetDefaultTemplateData(c, "Settings", "settings"), "layouts/main")

		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})
}
