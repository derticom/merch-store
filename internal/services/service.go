package services

import "github.com/derticom/merch-store/internal/repositories"

type Service struct {
	repo *repositories.Storage
}

func New(userRepo *repositories.Storage) Service {
	return Service{repo: userRepo}
}
