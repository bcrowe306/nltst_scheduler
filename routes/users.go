package routes

import (
	"github.com/gofiber/fiber/v3"

	"log"

	"github.com/bcrowe306/nltst_scheduler.git/models"
	"github.com/bcrowe306/nltst_scheduler.git/pages"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func CreateUsersRoutes(app *fiber.App, BaseRoute string) {
	app.Get(BaseRoute+"/new", Protected, func(c fiber.Ctx) error {
		// New user page
		data := GetDefaultTemplateData(c, "New User", BaseRoute)
		data["Title"] = "New User"
		err := c.Render("pages/users/new", data, "layouts/main")
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	app.Get(BaseRoute+"/:id", Protected, func(c fiber.Ctx) error {
		// User edit page
		db, ok := fiber.GetState[*mongo.Database](c.App().State(), "db")
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("Database not found in context")
		}

		userID := c.Params("id")
		user, err := models.FindUserByID(db, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching user")
		}

		data := GetDefaultTemplateData(c, "Edit User", BaseRoute)
		data["User"] = user

		err = c.Render("pages/users/edit", data, "layouts/main")

		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	app.Get(BaseRoute, Protected, func(c fiber.Ctx) error {
		db, ok := fiber.GetState[*mongo.Database](c.App().State(), "db")
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("Database not found in context")
		}
		data := GetDefaultTemplateData(c, "Users", BaseRoute)
		users, err := models.GetAllUsers(db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching users")
		}
		data["Users"] = users
		err = RenderHTMXPage(c, pages.UsersPage(data))
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template")
		}
		return nil
	})

	// New user handler
	app.Post(BaseRoute, Protected, func(c fiber.Ctx) error {
		db, ok := fiber.GetState[*mongo.Database](c.App().State(), "db")
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("Database not found in context")
		}

		var formData struct {
			Name        string `form:"name"`
			Email       string `form:"email"`
			PhoneNumber string `form:"phoneNumber"`
			Password    string `form:"password"`
		}

		if err := c.Bind().Form(&formData); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
		}

		if _, err := models.CreateUser(db, formData.Name, formData.Email, formData.Password, formData.PhoneNumber); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error creating user")
		} else {
			return c.Redirect().To("/users")
		}
	})

	// Update user handler
	app.Post(BaseRoute+"/:id", Protected, func(c fiber.Ctx) error {
		db, ok := fiber.GetState[*mongo.Database](c.App().State(), "db")
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("Database not found in context")
		}

		userID := c.Params("id")
		var updatedData models.User
		if err := c.Bind().Form(&updatedData); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
		}

		err := models.UpdateUser(db, userID, &updatedData)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error updating user")
		}

		return c.Redirect().To(BaseRoute)
	})

	// Delete user handler
	app.Get(BaseRoute+"/delete/:id", Protected, func(c fiber.Ctx) error {
		db, ok := fiber.GetState[*mongo.Database](c.App().State(), "db")
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("Database not found in context")
		}

		userID := c.Params("id")
		err := models.DeleteUser(db, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error deleting user")
		}

		return c.Redirect().To(BaseRoute)
	})

}
