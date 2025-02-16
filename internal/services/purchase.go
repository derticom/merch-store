package services

import (
	"context"
	"errors"

	"github.com/derticom/merch-store/internal/models"

	"github.com/google/uuid"
)

func (s *Service) BuyItem(ctx context.Context, userID uuid.UUID, itemName string) error {
	// Получаем данные пользователя и товара
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	item, err := s.repo.GetItemByName(ctx, itemName)
	if err != nil {
		return err
	}
	if item == nil {
		return errors.New("item not found")
	}

	// Проверяем, что у пользователя достаточно монет
	if user.Coins < item.Price {
		return errors.New("insufficient coins")
	}

	// Обновляем баланс пользователя
	if err := s.repo.UpdateUserCoins(ctx, userID, user.Coins-item.Price); err != nil {
		return err
	}

	// Создаем запись о покупке
	purchase := &models.Purchase{
		ID:     uuid.New(),
		UserID: userID,
		Item:   itemName,
	}

	return s.repo.CreatePurchase(ctx, purchase)
}

func (s *Service) GetPurchaseHistory(ctx context.Context, userID uuid.UUID) ([]models.Purchase, error) {
	return s.repo.GetPurchasesByUserID(ctx, userID)
}
