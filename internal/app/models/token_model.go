package models

import (
	"time"

	"github.com/google/uuid"
)

type TokenDetails struct {
	AccessToken      string
	RefreshToken     string
	AccessTokenUUID  string
	RefreshTokenUUID string
	AtExpiresAt      int64
	RtExpiresAt      int64
}

func (td *TokenDetails) CreateTokenDetails() {
	td.AtExpiresAt = time.Now().Add(time.Minute * time.Duration(15)).Unix()
	td.RtExpiresAt = time.Now().Add(time.Hour * time.Duration(24)).Unix()
	td.AccessTokenUUID = uuid.NewString()
	td.RefreshTokenUUID = uuid.NewString()
}
