package middleware

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/BravoAndres/fiber-api/pkg/auth"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	authManager, err := auth.NewManager(os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}

	errHandler := func(c *fiber.Ctx, err error) error {
		switch err.Error() {
		case "Token is expired":
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "Token expired",
			})
		case "missing or malformed JWT":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   "Missing or malformed JWT",
			})
		case "signature is invalid":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   "Bad Signature",
			})
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":    true,
				"msg":      "Invalid Token",
				"msgError": err.Error(),
			})
		}
	}

	successHandler := func(c *fiber.Ctx) error {
		return c.Next()
	}

	return func(c *fiber.Ctx) error {
		var auth string

		auth, err := extractTokenFromHeader(c)
		if err != nil {
			return errHandler(c, err)
		}

		token, err := authManager.Parse(auth)
		if err != nil {
			return errHandler(c, err)
		}

		// Store user information from token into context.
		c.Locals("userId", token)
		return successHandler(c)

	}
}

func extractTokenFromHeader(c *fiber.Ctx) (string, error) {
	bearerToken := c.Get("Authorization")

	// Get token from: "Bearer <token>"" format
	token := strings.Split(bearerToken, " ")
	if len(token) == 2 && strings.EqualFold(bearerToken[:len(token[0])], "Bearer") {
		return token[1], nil
	}

	return "", errors.New("missing or malformed JWT")
}
