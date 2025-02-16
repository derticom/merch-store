package repositories

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib" //revive:disable:blank-imports // import for side effect is necessary here.
	"github.com/pressly/goose/v3"
)

type Storage struct {
	db *sql.DB
}

func New(ctx context.Context, dsn string) (*Storage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to sql.Open: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to db.PingContext: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("failed to Close: %w", err)
	}

	return nil
}

func (s *Storage) Migrate(migrate string) (err error) {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to goose.SetDialect: %w", err)
	}

	if err := goose.Up(s.db, migrate); err != nil {
		return fmt.Errorf("failed to goose.Up: %w", err)
	}

	return nil
}
