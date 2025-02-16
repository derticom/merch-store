package services

import (
	"context"
	"errors"
	"testing"

	"github.com/derticom/merch-store/internal/models"
	"github.com/derticom/merch-store/internal/services/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService_BuyItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_services.NewMockRepository(ctrl)
	service := New(mockRepo)

	userID := uuid.New()
	itemName := "test-item"
	user := &models.User{
		ID:    userID,
		Coins: 100,
	}
	item := &models.Item{
		Name:  itemName,
		Price: 50,
	}

	tests := []struct {
		name        string
		setup       func()
		userID      uuid.UUID
		itemName    string
		expectedErr error
	}{
		{
			name: "successful purchase",
			setup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(user, nil)
				mockRepo.EXPECT().GetItemByName(gomock.Any(), itemName).Return(item, nil)
				mockRepo.EXPECT().UpdateUserCoins(gomock.Any(), userID, user.Coins-item.Price).Return(nil)
				mockRepo.EXPECT().CreatePurchase(gomock.Any(), gomock.Any()).Return(nil)
			},
			userID:      userID,
			itemName:    itemName,
			expectedErr: nil,
		},
		{
			name: "user not found",
			setup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(nil, nil)
			},
			userID:      userID,
			itemName:    itemName,
			expectedErr: errors.New("user not found"),
		},
		{
			name: "item not found",
			setup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(user, nil)
				mockRepo.EXPECT().GetItemByName(gomock.Any(), itemName).Return(nil, nil)
			},
			userID:      userID,
			itemName:    itemName,
			expectedErr: errors.New("item not found"),
		},
		{
			name: "insufficient coins",
			setup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(user, nil)
				mockRepo.EXPECT().GetItemByName(gomock.Any(), itemName).Return(&models.Item{Name: itemName, Price: 200}, nil)
			},
			userID:      userID,
			itemName:    itemName,
			expectedErr: errors.New("insufficient coins"),
		},
		{
			name: "error updating user coins",
			setup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(user, nil)
				mockRepo.EXPECT().GetItemByName(gomock.Any(), itemName).Return(item, nil)
				mockRepo.EXPECT().UpdateUserCoins(gomock.Any(), userID, user.Coins-item.Price).Return(errors.New("update error"))
			},
			userID:      userID,
			itemName:    itemName,
			expectedErr: errors.New("update error"),
		},
		{
			name: "error creating purchase",
			setup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(user, nil)
				mockRepo.EXPECT().GetItemByName(gomock.Any(), itemName).Return(item, nil)
				mockRepo.EXPECT().UpdateUserCoins(gomock.Any(), userID, user.Coins-item.Price).Return(nil)
				mockRepo.EXPECT().CreatePurchase(gomock.Any(), gomock.Any()).Return(errors.New("purchase error"))
			},
			userID:      userID,
			itemName:    itemName,
			expectedErr: errors.New("purchase error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := service.BuyItem(context.Background(), tt.userID, tt.itemName)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
