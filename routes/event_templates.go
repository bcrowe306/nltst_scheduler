package routes

import (
	"github.com/bcrowe306/nltst_scheduler.git/models"
	"github.com/gofiber/fiber/v3"

	"log"
)

func CreateEventTemplatesRoutes(app *fiber.App) {

	// Event Templates Index Route
	app.Get("/event_templates", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}
		data := GetDefaultTemplateData(c, "Event Templates", "event_templates")
		event_templates, err := models.GetAllEventTemplates(db)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching event templates")
		}
		data["EventTemplates"] = event_templates
		err = c.Render("pages/event_templates/index", data, "layouts/main")
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	// New Event Template Route
	app.Get("/event_templates/new", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}
		data := GetDefaultTemplateData(c, "New Event Template", "event_templates")
		teams, err := models.GetAllTeams(db)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching teams in event template creation")
		}

		data["Teams"] = teams
		err = c.Render("pages/event_templates/new", data, "layouts/main")
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	// Create Event Template Route
	app.Post("/event_templates", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		var eventTemplate models.EventTemplate
		if err := c.Bind().Form(&eventTemplate); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		if _, err := models.InsertEventTemplate(db, &eventTemplate); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error creating event template")
		}

		return c.Redirect().To("/event_templates")
	})

	// Edit Event Template Route
	app.Get("/event_templates/:id", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		eventTemplateID := c.Params("id")
		eventTemplate, err := models.GetEventTemplateByID(db, eventTemplateID)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching event template")
		}

		teams, err := models.GetAllTeams(db)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching teams in event template edit")
		}

		data := GetDefaultTemplateData(c, "Edit Event Template", "event_templates")
		data["EventTemplate"] = eventTemplate
		data["Teams"] = teams
		err = c.Render("pages/event_templates/edit", data, "layouts/main")
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	// Update Event Template Route
	app.Post("/event_templates/:id", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		eventTemplateID := c.Params("id")
		var eventTemplate models.EventTemplate
		if err := c.Bind().Form(&eventTemplate); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}
		eventTemplate.ID = eventTemplateID
		if _, err := models.UpdateEventTemplate(db, &eventTemplate); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error updating event template")
		}

		return c.Redirect().To("/event_templates")
	})

	// Delete Event Template Route
	app.Get("/event_templates/delete/:id", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		eventTemplateID := c.Params("id")
		if _, err := models.DeleteEventTemplate(db, eventTemplateID); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error deleting event template")
		}

		return c.Redirect().To("/event_templates")
	})

	// Add position to event template route
	app.Post("/event_templates/:id/positions", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		positionName := c.FormValue("position-name")
		positionDescription := c.FormValue("position-description")

		position := models.Position{
			Name:        positionName,
			Description: positionDescription,
		}

		eventTemplateID := c.Params("id")

		if _, err := models.AddPositionToEventTemplate(db, eventTemplateID, position); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error adding position to event template")
		}

		return c.Redirect().To("/event_templates/" + eventTemplateID)
	})

	// Remove position from event template route
	app.Get("/event_templates/:event_template_id/positions/:position_name/delete", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		eventTemplateID := c.Params("event_template_id")
		positionName := c.Params("position_name")

		if _, err := models.RemovePositionFromEventTemplate(db, eventTemplateID, positionName); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error removing position from event template")
		}

		return c.Redirect().To("/event_templates/" + eventTemplateID)
	})
}
