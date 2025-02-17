package models

import (
	"time"

	"github.com/google/uuid"
)

//nolint:tagliatelle // snake_case is allowed here.
type Purchase struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Item      string    `json:"item" db:"item"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
