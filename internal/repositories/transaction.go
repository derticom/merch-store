package repositories

import (
	"context"
	"errors"

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

func (s *Storage) SendCoins(ctx context.Context, fromUserID, toUserID uuid.UUID, amount int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var fromUserCoins int
	query := `SELECT coins FROM users WHERE id = $1`
	err = tx.QueryRowContext(ctx, query, fromUserID).Scan(&fromUserCoins)
	if err != nil {
		return err
	}

	if fromUserCoins < amount {
		return errors.New("insufficient coins")
	}

	query = `UPDATE users SET coins = coins - $1 WHERE id = $2`
	_, err = tx.ExecContext(ctx, query, amount, fromUserID)
	if err != nil {
		return err
	}

	query = `UPDATE users SET coins = coins + $1 WHERE id = $2`
	_, err = tx.ExecContext(ctx, query, amount, toUserID)
	if err != nil {
		return err
	}

	transaction := &models.Transaction{
		ID:       uuid.New(),
		FromUser: fromUserID,
		ToUser:   toUserID,
		Amount:   amount,
	}
	query = `INSERT INTO transactions (id, from_user, to_user, amount, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.ExecContext(ctx, query, transaction.ID, transaction.FromUser, transaction.ToUser,
		transaction.Amount, transaction.CreatedAt)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
