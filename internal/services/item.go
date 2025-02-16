package services

import (
	"context"

	"github.com/derticom/merch-store/internal/models"
)

func (s *Service) GetAllItems(ctx context.Context) ([]models.Item, error) {
	return s.repo.GetAllItems(ctx)
}

func (s *Service) GetItemByName(ctx context.Context, name string) (*models.Item, error) {
	return s.repo.GetItemByName(ctx, name)
}
