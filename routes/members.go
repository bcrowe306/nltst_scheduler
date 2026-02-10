package routes

import (
	"github.com/bcrowe306/nltst_scheduler.git/models"
	"github.com/bcrowe306/nltst_scheduler.git/pages"
	"github.com/gofiber/fiber/v3"

	"log"
)

func CreateMembersRoutes(app *fiber.App, BaseRoute string) {

	// Members Index
	app.Get(BaseRoute, Protected, func(c fiber.Ctx) error {
		data := GetDefaultTemplateData(c, "Members", BaseRoute)
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection not found")
		}

		members, err := models.GetAllMembers(db)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving members")
		}

		data["Members"] = members

		err = RenderHTMXPage(c, pages.MembersPage(data))
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	// New Member Form
	app.Get(BaseRoute+"/new", Protected, func(c fiber.Ctx) error {
		data := GetDefaultTemplateData(c, "New Member", BaseRoute)

		err := c.Render("pages/members/new", data, "layouts/main")
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	// Create Member
	app.Post(BaseRoute, Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection not found")
		}

		var new_member models.Member
		err = c.Bind().Form(&new_member)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
		}

		_, err = models.InsertMember(db, &new_member)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error creating member")
		}

		return c.Redirect().To(BaseRoute)
	})

	// Edit Member Form
	app.Get(BaseRoute+"/:id", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection not found")
		}

		memberID := c.Params("id")
		member, err := models.GetMemberByID(db, memberID)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving member")
		}

		data := GetDefaultTemplateData(c, "Edit Member", BaseRoute)
		data["Member"] = member

		err = RenderHTMXPage(c, pages.MembersEditPage(data))
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil

	})

	// Update Member
	app.Post(BaseRoute+"/:id", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection not found")
		}

		memberID := c.Params("id")
		var member = &models.Member{}

		err = c.Bind().Form(member)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
		}

		_, err = models.UpdateMember(db, memberID, member)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error updating member")
		}

		return c.Redirect().To(BaseRoute)
	})

	// Delete Member
	app.Get(BaseRoute+"/delete/:id", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection not found")
		}

		memberID := c.Params("id")
		_, err = models.DeleteMember(db, memberID)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error deleting member")
		}

		return c.Redirect().To(BaseRoute)
	})
}
