package routes

import (
	"github.com/bcrowe306/nltst_scheduler.git/models"
	"github.com/bcrowe306/nltst_scheduler.git/pages"
	"github.com/gofiber/fiber/v3"

	"log"
)

func CreateTeamsRoutes(app *fiber.App, BaseRoute string) {
	// Remove member from team route
	app.Get(BaseRoute+"/:id/remove_member/:member_id", Protected, func(c fiber.Ctx) error {
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

		return c.Redirect().To(BaseRoute + "/" + teamID)
	})

	// Teams Index Route
	app.Get(BaseRoute, Protected, func(c fiber.Ctx) error {
		data := GetDefaultTemplateData(c, "Teams", BaseRoute)
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

		err = RenderHTMXPage(c, pages.TeamsPage(data))
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	// New Team Form Route
	app.Get(BaseRoute+"/new", Protected, func(c fiber.Ctx) error {
		data := GetDefaultTemplateData(c, "New Team", BaseRoute)
		err := c.Render("pages/teams/new", data, "layouts/main")
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	// Create Team Route
	app.Post(BaseRoute, Protected, func(c fiber.Ctx) error {
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

		return c.Redirect().To(BaseRoute)
	})

	// View Team Route
	app.Get(BaseRoute+"/:id", Protected, func(c fiber.Ctx) error {
		data := GetDefaultTemplateData(c, "View Team", BaseRoute)
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
	app.Post(BaseRoute+"/:id/add_member", Protected, func(c fiber.Ctx) error {
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

		return c.Redirect().To(BaseRoute + "/" + teamID)
	})

	// Edit Team details route
	app.Post(BaseRoute+"/:id/edit", Protected, func(c fiber.Ctx) error {
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

		return c.Redirect().To(BaseRoute + "/" + teamID)
	})

	// Delete Team route
	app.Get(BaseRoute+"/:id/delete", Protected, func(c fiber.Ctx) error {
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

		return c.Redirect().To(BaseRoute)
	})
}
