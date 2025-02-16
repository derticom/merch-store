package services

import (
	"context"
	"errors"

	"github.com/derticom/merch-store/internal/models"

	"github.com/google/uuid"
)

func (s *Service) SendCoins(ctx context.Context, fromUserID, toUserID uuid.UUID, amount int) error {
	if fromUserID == toUserID {
		return errors.New("cannot send coins to yourself")
	}

	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	if err := s.repo.SendCoins(ctx, fromUserID, toUserID, amount); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetTransactionHistory(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	return s.repo.GetTransactionsByUserID(ctx, userID)
}
