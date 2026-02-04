package routes

import (
	"github.com/bcrowe306/nltst_scheduler.git/models"
	"github.com/gofiber/fiber/v3"

	"log"
)

func CreateTeamsRoutes(app *fiber.App) {
	// Remove member from team route
	app.Get("/teams/:id/remove_member/:member_id", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}
		teamID := c.Params("id")
		memberID := c.Params("member_id")

		_, err = models.RemoveMemberFromTeam(db, teamID, memberID)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error removing member from team")
		}

		return c.Redirect().To("/teams/" + teamID)
	})

	// Teams Index Route
	app.Get("/teams", Protected, func(c fiber.Ctx) error {
		data := GetDefaultTemplateData(c, "Teams", "teams")
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		teams, err := models.GetAllTeams(db)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching teams")
		}

		data["Teams"] = teams

		err = c.Render("pages/teams/index", data, "layouts/main")
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	// New Team Form Route
	app.Get("/teams/new", Protected, func(c fiber.Ctx) error {
		data := GetDefaultTemplateData(c, "New Team", "teams")
		err := c.Render("pages/teams/new", data, "layouts/main")
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	// Create Team Route
	app.Post("/teams", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		type NewTeamForm struct {
			Name        string `form:"name"`
			Description string `form:"description"`
		}
		var form NewTeamForm
		if err := c.Bind().Form(&form); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
		}

		team := &models.Team{
			Name:        form.Name,
			Description: form.Description,
			Members:     []string{},
		}

		_, err = models.InsertTeam(db, team)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error creating team")
		}

		return c.Redirect().To("/teams")
	})

	// View Team Route
	app.Get("/teams/:id", Protected, func(c fiber.Ctx) error {
		data := GetDefaultTemplateData(c, "View Team", "teams")
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		teamID := c.Params("id")
		team, err := models.GetTeamByID(db, teamID)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching team")
		}

		members, err := models.GetAllMembers(db)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching members")
		}

		data["Team"] = team
		data["Members"] = members

		err = c.Render("pages/teams/view", data, "layouts/main")
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	// Add member to team route
	app.Post("/teams/:id/add_member", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}
		teamID := c.Params("id")
		memberID := c.FormValue("memberID")

		_, err = models.AddMemberToTeam(db, teamID, memberID)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error adding member to team")
		}

		return c.Redirect().To("/teams/" + teamID)
	})

	// Edit Team details route
	app.Post("/teams/:id/edit", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}
		teamID := c.Params("id")

		type EditTeamForm struct {
			Name        string `form:"name"`
			Description string `form:"description"`
		}
		var form EditTeamForm
		if err := c.Bind().Form(&form); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
		}

		team := &models.Team{
			ID:          teamID,
			Name:        form.Name,
			Description: form.Description,
		}

		_, err = models.UpdateTeam(db, teamID, team)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error updating team")
		}

		return c.Redirect().To("/teams/" + teamID)
	})

	// Delete Team route
	app.Get("/teams/:id/delete", Protected, func(c fiber.Ctx) error {
		db, err := GetDatabaseFromContext(c)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}
		teamID := c.Params("id")

		_, err = models.DeleteTeam(db, teamID)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error deleting team")
		}

		return c.Redirect().To("/teams")
	})
}
