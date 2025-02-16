package models

type Item struct {
	Name  string `json:"name" db:"name"`
	Price int    `json:"price" db:"price"`
}
