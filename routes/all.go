package routes

import (
	"log"

	"github.com/a-h/templ"
	"github.com/bcrowe306/nltst_scheduler.git/pages"
	"github.com/gofiber/fiber/v3"
)

func isHTMXRequest(c fiber.Ctx) bool {
	if c.HasHeader("HX-Request") {
		hx_header := c.GetReqHeaders()["Hx-Request"]
		if hx_header[0] == "true" {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func Render(c fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	log.Print(isHTMXRequest(c))
	if isHTMXRequest(c) {
		// Handle HTMX request
		return component.Render(c.Context(), c.Response().BodyWriter())
	} else {
		return pages.Index(component).Render(c.Context(), c.Response().BodyWriter())
	}

}

func RenderHTMXPage(c fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	log.Print(isHTMXRequest(c))
	if isHTMXRequest(c) {
		// Handle HTMX request
		return component.Render(c.Context(), c.Response().BodyWriter())
	} else {
		return pages.Index(component).Render(c.Context(), c.Response().BodyWriter())
	}

}

func RenderFullPage(c fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

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
