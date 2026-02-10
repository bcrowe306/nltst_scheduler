package routes

import (
	"github.com/gofiber/fiber/v3"

	"log"
)

func CreateIntegrationsRoutes(app *fiber.App, BaseRoute string) {
	app.Get(BaseRoute, Protected, func(c fiber.Ctx) error {

		err := c.Render("pages/integrations/index", GetDefaultTemplateData(c, "Integrations", BaseRoute), "layouts/main")
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})
}
