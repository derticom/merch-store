package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/derticom/merch-store/internal/models"

	"github.com/google/uuid"
)

func (s *Storage) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, username, password, coins) VALUES ($1, $2, $3, $4)`
	_, err := s.db.ExecContext(ctx, query, user.ID, user.Username, user.Password, user.Coins)
	return err
}

func (s *Storage) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `SELECT id, username, password, coins FROM users WHERE id = $1`
	row := s.db.QueryRowContext(ctx, query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Coins)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *Storage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT id, username, password, coins FROM users WHERE username = $1`
	row := s.db.QueryRowContext(ctx, query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Coins)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *Storage) UpdateUserCoins(ctx context.Context, id uuid.UUID, coins int) error {
	query := `UPDATE users SET coins = $1 WHERE id = $2`
	_, err := s.db.ExecContext(ctx, query, coins, id)
	return err
}
