package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/BravoAndres/fiber-api/internal/app/controllers"
)

func UserRoutes(app *fiber.App) {
	route := app.Group("/api/v1")

	route.Get("/users", controllers.GetUsers)
	route.Post("/users", controllers.CreateUser)
	route.Get("/users/:id", controllers.GetUser)
}
