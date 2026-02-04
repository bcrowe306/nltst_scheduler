package routes

import (
	"github.com/bcrowe306/nltst_scheduler.git/models"
	"github.com/gofiber/fiber/v3"

	"log"
)

func CreateServicesRoutes(app *fiber.App) {

	// Services Index Route
	app.Get("/services", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}
		data := GetDefaultTemplateData(c, "Services")
		services, err := models.GetAllServices(db)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching services")
		}
		data["Services"] = services
		err = c.Render("pages/services/index", data, "layouts/main")

		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	// New Service Route
	app.Get("/services/new", Protected, func(c fiber.Ctx) error {
		data := GetDefaultTemplateData(c, "New Service")
		err := c.Render("pages/services/new", data, "layouts/main")
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	// Create Service Route
	app.Post("/services", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		var service models.Service
		if err := c.Bind().Form(&service); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		if _, err := models.InsertService(db, &service); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error creating service")
		}

		return c.Redirect().To("/services")
	})

	// Edit Service Route
	app.Get("/services/:id", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		serviceID := c.Params("id")
		service, err := models.GetServiceByID(db, serviceID)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching service")
		}

		data := GetDefaultTemplateData(c, "Edit Service")
		data["Service"] = service

		err = c.Render("pages/services/edit", data, "layouts/main")
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	// Update Service Route
	app.Post("/services/:id", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		serviceID := c.Params("id")
		var service models.Service
		if err := c.Bind().Form(&service); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}
		service.ID = serviceID

		if _, err := models.UpdateService(db, &service); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error updating service")
		}

		return c.Redirect().To("/services")
	})

	// Delete Service Route
	app.Get("/services/delete/:id", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		serviceID := c.Params("id")
		if _, err := models.DeleteService(db, serviceID); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error deleting service")
		}

		return c.Redirect().To("/services")
	})

	// Add position to service route
	app.Post("/services/:id/positions", Protected, func(c fiber.Ctx) error {
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

		serviceID := c.Params("id")

		if _, err := models.AddPositionToService(db, serviceID, position); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error adding position to service")
		}

		return c.Redirect().To("/services/" + serviceID)
	})

	// Remove position from service route
	app.Get("/services/:service_id/positions/:position_name/delete", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		serviceID := c.Params("service_id")
		positionName := c.Params("position_name")

		if _, err := models.RemovePositionFromService(db, serviceID, positionName); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error removing position from service")
		}

		return c.Redirect().To("/services/" + serviceID)
	})
}
