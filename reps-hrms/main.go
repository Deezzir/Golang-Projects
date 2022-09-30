package main

import (
	"context"
	"reps-hrms/pkg/models"
	"reps-hrms/pkg/routes"
	"reps-hrms/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ======================================
// Main
// ======================================
var server_port = "0.0.0.0:8080"

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer func() {
		if err := models.Mongo.Client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := fiber.New()
	routes.RegisterRoutes(app)

	if err := app.Listen(server_port); err != nil {
		utils.ErrLogger.Fatalln(err)
	}
}
