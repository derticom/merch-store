package services

import (
	"context"
	"errors"

	"github.com/derticom/merch-store/internal/models"

	"github.com/google/uuid"
)

func (s *Service) SendCoins(ctx context.Context, fromUserID, toUserID uuid.UUID, amount int) error {
	// Проверяем, что отправитель и получатель не совпадают
	if fromUserID == toUserID {
		return errors.New("cannot send coins to yourself")
	}

	// Проверяем, что сумма положительная
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	// Получаем данные отправителя и получателя
	fromUser, err := s.repo.GetUserByID(ctx, fromUserID)
	if err != nil {
		return err
	}
	if fromUser == nil {
		return errors.New("sender not found")
	}

	toUser, err := s.repo.GetUserByID(ctx, toUserID)
	if err != nil {
		return err
	}
	if toUser == nil {
		return errors.New("receiver not found")
	}

	// Проверяем, что у отправителя достаточно монет
	if fromUser.Coins < amount {
		return errors.New("insufficient coins")
	}

	// Обновляем балансы
	if err := s.repo.UpdateUserCoins(ctx, fromUserID, fromUser.Coins-amount); err != nil {
		return err
	}
	if err := s.repo.UpdateUserCoins(ctx, toUserID, toUser.Coins+amount); err != nil {
		return err
	}

	// Создаем транзакцию
	transaction := &models.Transaction{
		ID:       uuid.New(),
		FromUser: fromUserID,
		ToUser:   toUserID,
		Amount:   amount,
	}

	return s.repo.CreateTransaction(ctx, transaction)
}

func (s *Service) GetTransactionHistory(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	return s.repo.GetTransactionsByUserID(ctx, userID)
}
