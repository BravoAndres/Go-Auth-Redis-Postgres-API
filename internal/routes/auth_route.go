package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/BravoAndres/Go-Auth-Redis-Postgres-API/internal/app/controllers"
	"github.com/BravoAndres/Go-Auth-Redis-Postgres-API/internal/middleware"
)

func AuthRoutes(app *fiber.App) {
	route := app.Group("/api/v1")

	route.Post("/login", controllers.Login)
	route.Post("/token/new", controllers.RefreshToken)

	//Restricted Routes
	route.Post("/logout", middleware.AuthMiddleware(), controllers.Logout)
	route.Get("/protected", middleware.AuthMiddleware(), controllers.Protected)
}
