package routes

import (
	"github.com/gofiber/fiber/v3"

	"log"

	"github.com/bcrowe306/nltst_scheduler.git/pages"
)

func CreateSettingsRoutes(app *fiber.App, BaseRoute string) {
	app.Get(BaseRoute, Protected, func(c fiber.Ctx) error {
		data := GetDefaultTemplateData(c, "Settings", BaseRoute)
		err := RenderHTMXPage(c, pages.SettingsPage(data))

		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})
}
