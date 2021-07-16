package queries

import (
	"context"
	"errors"
	"time"

	"github.com/BravoAndres/fiber-api/internal/app/models"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type TokenQueries struct {
	*redis.Client
}

func (q *TokenQueries) SaveToken(td *models.TokenDetails, userId string) error {
	at := time.Unix(td.AtExpiresAt, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpiresAt, 0)
	now := time.Now()

	err := q.Set(ctx, td.AccessTokenUUID, userId, at.Sub(now)).Err()
	if err != nil {
		return err
	}

	err = q.Set(ctx, td.RefreshTokenUUID, userId, rt.Sub(now)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (q *TokenQueries) FetchToken(token string) (string, error) {
	userId, err := q.Get(ctx, token).Result()
	if err == redis.Nil {
		return "", errors.New("key does not exist")
	} else if err != nil {
		return "", err
	}

	return userId, nil
}

func (q *TokenQueries) DeleteToken(uuid string) (int64, error) {
	deleted, err := q.Del(ctx, uuid).Result()
	if err != nil {
		return 0, err
	}

	return deleted, nil
}
