package controllers

import (
	"os"
	"strconv"

	"github.com/BravoAndres/Go-Auth-Redis-Postgres-API/internal/app/models"
	"github.com/BravoAndres/Go-Auth-Redis-Postgres-API/pkg/auth"
	"github.com/BravoAndres/Go-Auth-Redis-Postgres-API/pkg/database"
	"github.com/BravoAndres/Go-Auth-Redis-Postgres-API/pkg/hasher"
	"github.com/gofiber/fiber/v2"
)

type Token struct {
	RefreshToken string `json:"refresh_token"`
}

func Register(c *fiber.Ctx) error {
	reqUser := &models.User{}
	if err := c.BodyParser(reqUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
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

	user, err := db.GetUserByEmail(reqUser.Email)
	if err != nil && (models.User{} == user) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "User already exists",
		})
	}

	err = db.CreateUser(reqUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Error registering",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"user":  reqUser.Email,
	})
}

func Login(c *fiber.Ctx) error {
	reqUser := &models.User{}
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

	td := &models.TokenDetails{}
	td.CreateTokenDetails()
	tokenPair, err := tokenManager.GenerateTokenPairWithClaims(strconv.Itoa(user.ID),
		td.AccessTokenUUID,
		td.RefreshTokenUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	td.AccessToken = tokenPair["access_token"]
	td.RefreshToken = tokenPair["refresh_token"]

	err = db.SaveToken(td, strconv.Itoa(user.ID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}

	return c.JSON(fiber.Map{
		"access_token":  tokenPair["access_token"],
		"token_type":    "bearer",
		"expires_in":    td.AtExpiresAt,
		"refresh_token": tokenPair["refresh_token"],
	})
}

func Logout(c *fiber.Ctx) error {
	db, err := database.ConnectDB()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	token, ok := c.Locals("refresh_uuid").(string)
	if ok {
		deleted, err := db.DeleteToken(token)
		if err != nil && deleted == 0 {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
	}
	// Clean c.Locals
	c.Locals("userId", "")
	c.Locals("refresh_uuid", "")

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "User successfully logged out",
	})

}

func RefreshToken(c *fiber.Ctx) error {
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

	claims, err := tokenManager.Parse(reqToken.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":  true,
			"msg":    "Unauthorized",
			"msgErr": err.Error(),
		})
	}

	_, ok := claims["access_uuid"].(string)
	if ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Please use Refresh Token to get new access token.",
		})
	}

	db, err := database.ConnectDB()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	refreshUuid, ok := claims["refresh_uuid"].(string)
	if ok {
		// Check if refresh token exists in database
		// if not, should ask for new token logging in
		_, err := db.FetchToken(refreshUuid)
		if err != nil {
			if err.Error() == "key does not exist" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": true,
					"msg":   "Token is not valid",
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":    true,
				"msg":      "Unexpected Error",
				"msgError": err.Error(),
			})
		}
		deleted, err := db.DeleteToken(refreshUuid)
		if err != nil && deleted == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
	}

	userId, ok := claims["sub"].(string)
	if !ok {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	td := &models.TokenDetails{}
	td.CreateTokenDetails()
	tokenPair, err := tokenManager.GenerateTokenPairWithClaims(userId,
		td.AccessTokenUUID,
		td.RefreshTokenUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	td.AccessToken = tokenPair["access_token"]
	td.RefreshToken = tokenPair["refresh_token"]

	err = db.SaveToken(td, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err,
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
