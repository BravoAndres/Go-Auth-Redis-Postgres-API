package queries

import (
	"os"

	"github.com/jmoiron/sqlx"

	"github.com/BravoAndres/fiber-api/internal/app/models"
	"github.com/BravoAndres/fiber-api/pkg/hasher"
	"github.com/BravoAndres/fiber-api/pkg/logger"
)

type UserQueries struct {
	*sqlx.DB
}

// GetUsers method for getting all users.
func (q *UserQueries) GetUsers() ([]models.User, error) {
	// Define users variable.
	users := []models.User{}

	// Define query string.
	query := `SELECT * FROM users`

	// Send query to database.
	err := q.Select(&users, query)
	if err != nil {
		// Return empty object and error.
		return users, err
	}

	// Return query result.
	return users, nil
}

// GetUser method for getting one user by given ID.
func (q *UserQueries) GetUserById(id int) (models.User, error) {
	// Define user variable.
	user := models.User{}

	// Define query string.
	query := `SELECT * FROM users WHERE id = $1`

	// Send query to database.
	err := q.Get(&user, query, id)
	if err != nil {
		// Return empty object and error.
		return user, err
	}

	// Return query result.
	return user, nil
}

func (q *UserQueries) GetUserByCredentials(email string, password string) (models.User, error) {
	user := models.User{}

	query := `SELECT * FROM users WHERE email = $1 and password = $2`

	err := q.Get(&user, query, email, password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (q *UserQueries) CreateUser(u *models.User) error {
	hasher := hasher.NewSHA1Hasher(os.Getenv("PASSWORD_HASH_SALT"))
	hashedPassword, err := hasher.Hash(u.Password)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	query := `INSERT INTO users (email, password) VALUES ($1, $2)`

	_, err = q.Exec(query, u.Email, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}
