package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Username string    `json:"username" db:"username"`
	Password string    `json:"-" db:"password"`
	Coins    int       `json:"coins" db:"coins"`
}
