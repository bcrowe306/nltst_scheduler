package routes

import (
	"log"

	"github.com/a-h/templ"
	"github.com/bcrowe306/nltst_scheduler.git/pages"
	"github.com/gofiber/fiber/v3"
)

func isHTMXRequest(c fiber.Ctx) bool {
	result := false
	if c.HasHeader("HX-Request") {
		hx_header := c.GetReqHeaders()["Hx-Request"]
		if hx_header[0] == "true" {
			result = true
		} else {
			result = false
		}
	} else {
		result = false
	}
	log.Printf("isHTMXRequest: %v", result)
	return result
}

func Render(c fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")

	if isHTMXRequest(c) {
		// Handle HTMX request
		return component.Render(c.Context(), c.Response().BodyWriter())
	} else {
		return pages.Index(component).Render(c.Context(), c.Response().BodyWriter())
	}

}

func RenderHTMXPage(c fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
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
	CreateDashboardRoutes(app, "/dashboard")
	CreateMembersRoutes(app, "/members")
	CreateTeamsRoutes(app, "/teams")
	CreateEventTemplatesRoutes(app, "/event_templates")
	CreateScheduleRoutes(app, "/schedule")
	CreateIntegrationsRoutes(app, "/integrations")
	CreateUsersRoutes(app, "/users")
	CreateSettingsRoutes(app, "/settings")
	CreateAuthRoutes(app, "/auth")
}
