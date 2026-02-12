package routes

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/bcrowe306/nltst_scheduler.git/models"
	"github.com/bcrowe306/nltst_scheduler.git/pages"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var RoutePathMap = map[string]string{
	"event_templates": "Templates",
	"dashboard":       "Dashboard",
	"teams":           "Teams",
	"members":         "Members",
	"settings":        "Settings",
	"users":           "Users",
	"schedule":        "Schedule",
	"new":             "New",
	"edit":            "Edit",
}

func BreadcrumbMiddleware(c fiber.Ctx) error {
	// Middleware logic here
	breadcrumbs := GetRoutePathList(c)
	c.SetContext(context.WithValue(c.Context(), "Breadcrumbs", breadcrumbs))
	return c.Next()
}

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
	// log.Printf("isHTMXRequest: %v", result)
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

func Protected(c fiber.Ctx) error {
	sess := session.FromContext(c)
	if sess == nil {
		log.Print("Session is Nil")
		return c.Redirect().To("/login")
	}

	// Check if user is authenticated
	if sess.Get("authenticated") != true {
		log.Print("Not Authenticated")
		return c.Redirect().To("/login")
	}

	return c.Next()
}

func GetUserFromSession(c fiber.Ctx) (*models.User, error) {
	sess := session.FromContext(c)
	if sess == nil {
		return nil, fiber.ErrUnauthorized
	}

	userID := sess.Get("user_id")
	if userID == nil {
		return nil, fiber.ErrUnauthorized
	}

	db, ok := fiber.GetState[*mongo.Database](c.App().State(), "db")
	if !ok {
		return nil, fiber.ErrInternalServerError
	}

	user, err := models.FindUserByID(db, userID.(string))
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	return user, nil
}

func GetRoutePathList(c fiber.Ctx) []models.Breadcrumb {
	path := c.Path()

	// Split the path into segments by /
	var segments []string
	var breadcrumbs []models.Breadcrumb
	for _, segment := range strings.Split(path, "/") {
		if segment != "" {

			if RoutePathMap[segment] != "" {
				segment = RoutePathMap[segment]
			}
			segments = append(segments, segment)
			breadcrumbs = append(breadcrumbs, models.Breadcrumb{
				Label: segment,
				Url:   "/" + strings.Join(segments, "/"),
			})
		}
	}
	return breadcrumbs
}

func GetDefaultTemplateData(c fiber.Ctx, title string, sidebar_nav string) fiber.Map {
	user, err := GetUserFromSession(c)
	if err != nil {
		return fiber.Map{}
	}

	c.SetContext(context.WithValue(c.Context(), "SidebarNav", sidebar_nav))
	return fiber.Map{
		"Title":       title,
		"TimeNow":     time.Now(),
		"Breadcrumbs": GetRoutePathList(c),
		"User":        user,
		"SidebarNav":  sidebar_nav,
	}
}

func GetDatabaseFromContext(c fiber.Ctx) (*mongo.Database, error) {
	db, ok := fiber.GetState[*mongo.Database](c.App().State(), "db")
	if !ok {
		return nil, fiber.ErrInternalServerError
	}
	return db, nil
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
