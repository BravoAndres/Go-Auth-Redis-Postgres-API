package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/BravoAndres/Go-Auth-Redis-Postgres-API/internal/app/controllers"
	"github.com/BravoAndres/Go-Auth-Redis-Postgres-API/internal/middleware"
)

func UserRoutes(app *fiber.App) {
	route := app.Group("/api/v1")

	route.Get("/users", middleware.AuthMiddleware(), controllers.GetUsers)
	route.Post("/users", middleware.AuthMiddleware(), controllers.CreateUser)
	route.Get("/users/:id", middleware.AuthMiddleware(), controllers.GetUser)
}
