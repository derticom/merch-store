package handlers

import (
	"context"

	"github.com/derticom/merch-store/internal/models"

	"github.com/google/uuid"
)

//go:generate go run github.com/golang/mock/mockgen  -destination=mocks/mock_service.go . Service
type Service interface {
	RegisterUser(ctx context.Context, username, password string) (*models.User, error)
	AuthenticateUser(ctx context.Context, username, password string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateUserCoins(ctx context.Context, id uuid.UUID, coins int) error
	GetAllItems(ctx context.Context) ([]models.Item, error)
	GetItemByName(ctx context.Context, name string) (*models.Item, error)
	BuyItem(ctx context.Context, userID uuid.UUID, itemName string) error
	GetPurchaseHistory(ctx context.Context, userID uuid.UUID) ([]models.Purchase, error)
	SendCoins(ctx context.Context, fromUserID, toUserID uuid.UUID, amount int) error
	GetTransactionHistory(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error)
}

type Handler struct {
	service Service
}

func New(service Service) *Handler {
	return &Handler{service: service}
}
