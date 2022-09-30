package routes

import (
	"reps-hrms/pkg/models"

	"github.com/gofiber/fiber/v2"
)

var RegisterRoutes = func(app *fiber.App) {
	app.Get("/repository", models.GetRepositories)
	app.Post("/repository", models.CreateRepository)
	app.Post("/repository/upload", models.UploadRepositories)
	app.Get("/repository/:id", models.GetRepositoryByID)
	app.Put("/repository/:id", models.UpdateRepositoryByID)
	app.Delete("/repository/:id", models.DeleteRepository)
}
