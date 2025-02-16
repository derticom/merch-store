package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/derticom/merch-store/internal/models"
)

type ItemRepository interface {
	GetAllItems(ctx context.Context) ([]models.Item, error)
	GetItemByName(ctx context.Context, name string) (*models.Item, error)
}

func (s *Storage) GetAllItems(ctx context.Context) ([]models.Item, error) {
	query := `SELECT name, price FROM items`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.Name, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Storage) GetItemByName(ctx context.Context, name string) (*models.Item, error) {
	query := `SELECT name, price FROM items WHERE name = $1`
	row := s.db.QueryRowContext(ctx, query, name)

	var item models.Item
	err := row.Scan(&item.Name, &item.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}
