package routes

import (
	"github.com/gofiber/fiber/v3"
)

func CreateAllRoutes(app *fiber.App) {
	CreateDashboardRoutes(app)
	CreateMembersRoutes(app)
	CreateTeamsRoutes(app)
	CreateServicesRoutes(app)
	CreateSchedulesRoutes(app)
	CreateIntegrationsRoutes(app)
	CreateUsersRoutes(app)
}
