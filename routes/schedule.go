package routes

import (
	"time"

	"github.com/bcrowe306/nltst_scheduler.git/models"
	"github.com/bcrowe306/nltst_scheduler.git/pages"
	"github.com/gofiber/fiber/v3"

	"log"
)

func CreateScheduleRoutes(app *fiber.App, BaseRoute string) {

	// Schedule main page
	app.Get(BaseRoute, Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		events, err := models.GetEventsWithMemberDetails(db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching events")
		}

		event_templates, err := models.GetAllEventTemplates(db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching event templates")
		}

		data := GetDefaultTemplateData(c, "Schedule", BaseRoute)
		data["Events"] = events
		data["EventTemplates"] = event_templates

		err = RenderHTMXPage(c, pages.SchedulePage())
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	// Remove Position from Event
	app.Get(BaseRoute+"/:event_id/positions/delete/:position_name", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}
		eventID := c.Params("event_id")
		positionName := c.Params("position_name")

		_, err = models.RemovePosition(db, eventID, positionName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error removing position from event")
		}

		return c.Redirect().To(BaseRoute + "/" + eventID)
	})

	// Delete Event
	app.Get(BaseRoute+"/delete/:event_id", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		eventID := c.Params("event_id")
		_, err = models.DeleteEvent(db, eventID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error deleting event")
		}

		return c.Redirect().To(BaseRoute)
	})

	// Edit Event in schedule
	app.Get(BaseRoute+"/:event_id", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		eventID := c.Params("event_id")
		event, err := models.GetEventByID(db, eventID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching event")
		}

		var teamMembers []models.Member
		if event.TeamID == "" {
			teamMembers, err = models.GetAllMembers(db)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Error fetching all members")
			}
		} else {
			teamMembers, err = models.GetTeamMembers(db, event.TeamID)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Error fetching team members")
			}
		}

		data := GetDefaultTemplateData(c, "Edit Event", BaseRoute)
		data["Event"] = event
		data["TeamMembers"] = teamMembers

		err = c.Render("pages/schedule/edit", data, "layouts/main")
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	// New event in schedule
	app.Post(BaseRoute, Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		eventTemplateID := c.FormValue("eventTemplateID")
		eventDate := c.FormValue("eventDate")
		var eventTemplate *models.EventTemplate
		if eventTemplateID != "" {
			eventTemplate, err = models.GetEventTemplateByID(db, eventTemplateID)
			if err != nil {
				log.Print(err)
				return c.Status(fiber.StatusInternalServerError).SendString("Error fetching event template")
			}
		}
		parsed_date, err := time.Parse("2006-01-02", eventDate)
		if err != nil {
			log.Print(err)
			parsed_date = time.Now()
		}
		temp_event := models.Event{
			Name: "New Event",
			Date: parsed_date,
		}

		// If an event template was selected, copy its details to the new event
		if eventTemplate != nil {
			temp_event.StartTime = eventTemplate.StartTime
			temp_event.EndTime = eventTemplate.EndTime
			temp_event.Name = eventTemplate.Name
			temp_event.Description = eventTemplate.Description
			temp_event.TeamID = eventTemplate.TeamID

		}

		// Insert the new event into the database
		res, err := models.InsertEvent(db, &temp_event)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error creating new event")
		}
		new_event_id := res.InsertedID.(string)
		new_event, err := models.GetEventByID(db, new_event_id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching new event")
		}

		// If an event template was selected, add its positions to the new event
		if eventTemplate != nil {
			for _, pos := range eventTemplate.Positions {
				_, err = models.AddPosition(db, new_event.ID, models.PositionAssignment{
					PositionName: pos.Name,
					Description:  pos.Description,
				})
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).SendString("Error adding position to new event")
				}
			}
		}

		return c.Redirect().To(BaseRoute + "/" + new_event.ID)

	})

	// Update Event in schedule
	app.Post(BaseRoute+"/:event_id", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		eventID := c.Params("event_id")
		event, err := models.GetEventByID(db, eventID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching event")
		}

		event.Name = c.FormValue("name")
		event.Description = c.FormValue("description")
		// Date from date string
		date_str := c.FormValue("date")
		date, err := time.Parse("2006-01-02", date_str)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error parsing date")
		}
		event.Date = date
		event.StartTime = c.FormValue("startTime")
		event.EndTime = c.FormValue("endTime")

		_, err = models.UpdateEvent(db, event)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error updating event")
		}

		return c.Redirect().To(BaseRoute + "/" + eventID)
	})

	// Add Position to Event
	app.Post(BaseRoute+"/:event_id/positions", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}
		eventID := c.Params("event_id")
		positionName := c.FormValue("position_name")

		position := models.PositionAssignment{
			PositionName: positionName,
		}

		_, err = models.AddPosition(db, eventID, position)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error adding position to event")
		}

		return c.Redirect().To(BaseRoute + "/" + eventID)
	})

	// Assign Position to Member
	app.Post(BaseRoute+"/:event_id/positions/assign", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}
		eventID := c.Params("event_id")
		positionID := c.FormValue("positionID")
		memberID := c.FormValue("member_id")

		_, err = models.AssignPositionToMember(db, eventID, positionID, memberID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error assigning position to member")
		}

		return c.Redirect().To(BaseRoute + "/" + eventID)
	})

	// Unassign Position from Member
	app.Post(BaseRoute+"/:event_id/positions/unassign", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}
		eventID := c.Params("event_id")
		positionID := c.FormValue("positionID")

		_, err = models.UnassignPositionFromMember(db, eventID, positionID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error unassigning position from member")
		}

		return c.Redirect().To(BaseRoute + "/" + eventID)
	})

}
