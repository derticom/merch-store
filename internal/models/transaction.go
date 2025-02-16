package models

import (
	"github.com/google/uuid"
)

//nolint:tagliatelle // snake_case is allowed here.
type Transaction struct {
	ID        uuid.UUID `json:"id" db:"id"`
	FromUser  uuid.UUID `json:"from_user" db:"from_user"`
	ToUser    uuid.UUID `json:"to_user" db:"to_user"`
	Amount    int       `json:"amount" db:"amount"`
	CreatedAt int       `json:"created_at" db:"created_at"`
}
