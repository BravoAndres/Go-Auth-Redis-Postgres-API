package controllers

import (
	"os"
	"strconv"

	"github.com/BravoAndres/fiber-api/pkg/auth"
	"github.com/BravoAndres/fiber-api/pkg/database"
	"github.com/BravoAndres/fiber-api/pkg/hasher"
	"github.com/gofiber/fiber/v2"
)

type ReqUser struct {
	Email    string `db:"email" json:"email" validate:"required"`
	Password string `db:"password" json:"password" validate:"required,lte=255"`
}

type Token struct {
	RefreshToken string `json:"refresh_token"`
}

func Login(c *fiber.Ctx) error {
	reqUser := &ReqUser{}

	if err := c.BodyParser(reqUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	hashedPassword, err := hasher.NewSHA1Hasher(os.Getenv("PASSWORD_HASH_SALT")).Hash(reqUser.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	db, err := database.ConnectDB()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	user, err := db.GetUserByCredentials(reqUser.Email, hashedPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "User does not exists",
		})
	}

	tokenManager, err := auth.NewManager(os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	tokenPair, err := tokenManager.GenerateTokenPair(strconv.Itoa(user.ID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":         false,
		"msg":           nil,
		"access_token":  tokenPair["access_token"],
		"refresh_token": tokenPair["refresh_token"],
	})
}

func GetNewAccessToken(c *fiber.Ctx) error {
	reqToken := &Token{}

	if err := c.BodyParser(reqToken); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	tokenManager, err := auth.NewManager(os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	userId, err := tokenManager.Parse(reqToken.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":  true,
			"msg":    "Unauthorized",
			"msgErr": err.Error(),
		})
	}

	tokenPair, err := tokenManager.GenerateTokenPair(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":         false,
		"msg":           nil,
		"access_token":  tokenPair["access_token"],
		"refresh_token": tokenPair["refresh_token"],
	})
}

// Protected content controller
func Protected(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error":   false,
		"msg":     nil,
		"user_id": c.Locals("userId"),
	})
}
