package database

import (
	"github.com/BravoAndres/fiber-api/internal/app/queries"
)

type Queries struct {
	*queries.UserQueries
}

func ConnectDB() (*Queries, error) {
	db, err := PostgreSQLConnection()
	if err != nil {
		return nil, err
	}

	return &Queries{
		UserQueries: &queries.UserQueries{DB: db},
	}, nil
}
