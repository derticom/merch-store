package repositories

import (
	"context"

	"github.com/derticom/merch-store/internal/models"

	"github.com/google/uuid"
)

func (s *Storage) CreatePurchase(ctx context.Context, purchase *models.Purchase) error {
	query := `INSERT INTO purchase (id, user_id, item) VALUES ($1, $2, $3)`
	_, err := s.db.ExecContext(ctx, query, purchase.ID, purchase.UserID, purchase.Item)
	return err
}

func (s *Storage) GetPurchasesByUserID(ctx context.Context, userID uuid.UUID) ([]models.Purchase, error) {
	query := `SELECT id, user_id, item, created_at FROM purchase WHERE user_id = $1`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []models.Purchase
	for rows.Next() {
		var purchase models.Purchase
		if err := rows.Scan(&purchase.ID, &purchase.UserID, &purchase.Item, &purchase.CreatedAt); err != nil {
			return nil, err
		}
		purchases = append(purchases, purchase)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return purchases, nil
}
