package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID                int            `db:"id" json:"id,omitempty" validate:"required"`
	CreatedAt         time.Time      `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt         sql.NullTime   `db:"updated_at" json:"updated_at,omitempty"`
	Email             string         `db:"email" json:"email" validate:"required"`
	Password          string         `db:"password" json:"password" validate:"required,lte=255"`
	Role              string         `db:"role" json:"role,omitempty" validate:"required,lte=255"`
	VerificationToken sql.NullString `db:"verification_token" json:"verification_token,omitempty"`
	UserStatus        int            `db:"user_status" json:"user_status,omitempty" validate:"required,len=1"`
}
