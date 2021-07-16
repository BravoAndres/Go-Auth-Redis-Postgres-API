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

func (q *UserQueries) GetUsers() ([]models.User, error) {
	users := []models.User{}

	query := `SELECT * FROM users`

	err := q.Select(&users, query)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (q *UserQueries) GetUserById(id int) (models.User, error) {
	user := models.User{}

	query := `SELECT * FROM users WHERE id = $1`

	err := q.Get(&user, query, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (q *UserQueries) GetUserByEmail(email string) (models.User, error) {
	user := models.User{}

	query := `SELECT * FROM users WHERE email = $1`

	err := q.Get(&user, query, email)
	if err != nil {
		return user, err
	}

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
	hashedPassword, err := hasher.NewSHA1Hasher(os.Getenv("PASSWORD_HASH_SALT")).Hash(u.Password)
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
