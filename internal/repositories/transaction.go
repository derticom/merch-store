package repositories

import (
	"context"

	"github.com/derticom/merch-store/internal/models"

	"github.com/google/uuid"
)

func (s *Storage) CreateTransaction(ctx context.Context, transaction *models.Transaction) error {
	query := `INSERT INTO transactions (id, from_user, to_user, amount) VALUES ($1, $2, $3, $4)`
	_, err := s.db.ExecContext(ctx, query, transaction.ID, transaction.FromUser, transaction.ToUser, transaction.Amount)
	return err
}

func (s *Storage) GetTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	query := `SELECT id, from_user, to_user, amount, created_at FROM transactions WHERE from_user = $1 OR to_user = $1`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		if err = rows.Scan(
			&transaction.ID,
			&transaction.FromUser,
			&transaction.ToUser,
			&transaction.Amount,
			&transaction.CreatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
