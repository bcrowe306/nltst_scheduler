package routes

import (
	"github.com/bcrowe306/nltst_scheduler.git/models"
	"github.com/bcrowe306/nltst_scheduler.git/pages"
	"github.com/gofiber/fiber/v3"
)

func CreateDashboardRoutes(app *fiber.App, BaseRoute string) {
	app.Get("/", Protected, func(c fiber.Ctx) error {

		return c.Redirect().To(BaseRoute)
	})
	app.Get(BaseRoute, Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		data := GetDefaultTemplateData(c, "Dashboard", BaseRoute)
		eventsByPosition, err := models.GetEventsGroupedByPositionName(db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching events")
		}

		events, err := models.GetEventsWithMemberDetails(db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching events")
		}

		data["Events"] = events

		data["Positions"] = eventsByPosition
		RenderHTMXPage(c, pages.DashboardPage(data))

		return nil
	})
}
