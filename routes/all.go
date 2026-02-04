package routes

import (
	"github.com/gofiber/fiber/v3"
)

func CreateAllRoutes(app *fiber.App) {
	CreateDashboardRoutes(app)
	CreateMembersRoutes(app)
	CreateTeamsRoutes(app)
	CreateEventTemplatesRoutes(app)
	CreateScheduleRoutes(app)
	CreateIntegrationsRoutes(app)
	CreateUsersRoutes(app)
	CreateSettingsRoutes(app)
	CreateAuthRoutes(app)
}
