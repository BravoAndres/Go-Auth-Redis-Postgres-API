package database

import (
	"github.com/BravoAndres/Go-Auth-Redis-Postgres-API/internal/app/queries"
)

type Queries struct {
	*queries.UserQueries
	*queries.TokenQueries
}

func ConnectDB() (*Queries, error) {
	pgdb, err := PostgreSQLConnection()
	if err != nil {
		return nil, err
	}

	rdb, err := NewRedisClient()
	if err != nil {
		return nil, err
	}

	return &Queries{
		UserQueries:  &queries.UserQueries{DB: pgdb},
		TokenQueries: &queries.TokenQueries{Client: rdb},
	}, nil
}
