package routes

import (
	"github.com/gofiber/fiber/v2"

	swagger "github.com/arsmn/fiber-swagger/v2" // fiber-swagger middleware
)

// SwaggerRoute func for describe group of API Docs routes.
func SwaggerRoute(app *fiber.App) {
	// Create routes group.
	route := app.Group("/swagger")

	// Routes for GET method:
	route.Get("*", swagger.Handler) // get one user by ID
}
