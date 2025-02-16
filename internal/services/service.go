package services

import (
	"context"

	"github.com/derticom/merch-store/internal/models"

	"github.com/google/uuid"
)

//go:generate go run github.com/golang/mock/mockgen  -destination=mocks/mock_repository.go . Repository
type Repository interface {
	GetAllItems(ctx context.Context) ([]models.Item, error)
	GetItemByName(ctx context.Context, name string) (*models.Item, error)
	CreatePurchase(ctx context.Context, purchase *models.Purchase) error
	GetPurchasesByUserID(ctx context.Context, userID uuid.UUID) ([]models.Purchase, error)
	CreateTransaction(ctx context.Context, transaction *models.Transaction) error
	GetTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error)
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUserCoins(ctx context.Context, id uuid.UUID, coins int) error
	SendCoins(ctx context.Context, fromUserID, toUserID uuid.UUID, amount int) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}
