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

func TestService_SendCoins(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_services.NewMockRepository(ctrl)
	service := New(mockRepo)

	fromUserID := uuid.New()
	toUserID := uuid.New()
	amount := 50

	fromUser := &models.User{
		ID:    fromUserID,
		Coins: 100,
	}
	toUser := &models.User{
		ID:    toUserID,
		Coins: 50,
	}

	tests := []struct {
		name        string
		setup       func()
		fromUserID  uuid.UUID
		toUserID    uuid.UUID
		amount      int
		expectedErr error
	}{
		{
			name: "successful coin transfer",
			setup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), fromUserID).Return(fromUser, nil)
				mockRepo.EXPECT().GetUserByID(gomock.Any(), toUserID).Return(toUser, nil)
				mockRepo.EXPECT().UpdateUserCoins(gomock.Any(), fromUserID, fromUser.Coins-amount).Return(nil)
				mockRepo.EXPECT().UpdateUserCoins(gomock.Any(), toUserID, toUser.Coins+amount).Return(nil)
				mockRepo.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(nil)
			},
			fromUserID:  fromUserID,
			toUserID:    toUserID,
			amount:      amount,
			expectedErr: nil,
		},
		{
			name: "sender and receiver are the same",
			setup: func() {
				// Никаких вызовов репозитория, так как проверка происходит до обращения к репозиторию
			},
			fromUserID:  fromUserID,
			toUserID:    fromUserID, // Отправитель и получатель совпадают
			amount:      amount,
			expectedErr: errors.New("cannot send coins to yourself"),
		},
		{
			name: "non-positive amount",
			setup: func() {
				// Никаких вызовов репозитория, так как проверка происходит до обращения к репозиторию
			},
			fromUserID:  fromUserID,
			toUserID:    toUserID,
			amount:      0, // Нулевая сумма
			expectedErr: errors.New("amount must be positive"),
		},
		{
			name: "sender not found",
			setup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), fromUserID).Return(nil, nil)
			},
			fromUserID:  fromUserID,
			toUserID:    toUserID,
			amount:      amount,
			expectedErr: errors.New("sender not found"),
		},
		{
			name: "receiver not found",
			setup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), fromUserID).Return(fromUser, nil)
				mockRepo.EXPECT().GetUserByID(gomock.Any(), toUserID).Return(nil, nil)
			},
			fromUserID:  fromUserID,
			toUserID:    toUserID,
			amount:      amount,
			expectedErr: errors.New("receiver not found"),
		},
		{
			name: "insufficient coins",
			setup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), fromUserID).Return(fromUser, nil)
				mockRepo.EXPECT().GetUserByID(gomock.Any(), toUserID).Return(toUser, nil)
			},
			fromUserID:  fromUserID,
			toUserID:    toUserID,
			amount:      150, // Сумма больше, чем есть у отправителя
			expectedErr: errors.New("insufficient coins"),
		},
		{
			name: "error updating sender's coins",
			setup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), fromUserID).Return(fromUser, nil)
				mockRepo.EXPECT().GetUserByID(gomock.Any(), toUserID).Return(toUser, nil)
				mockRepo.EXPECT().UpdateUserCoins(gomock.Any(), fromUserID, fromUser.Coins-amount).Return(
					errors.New("update error"))
			},
			fromUserID:  fromUserID,
			toUserID:    toUserID,
			amount:      amount,
			expectedErr: errors.New("update error"),
		},
		{
			name: "error updating receiver's coins",
			setup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), fromUserID).Return(fromUser, nil)
				mockRepo.EXPECT().GetUserByID(gomock.Any(), toUserID).Return(toUser, nil)
				mockRepo.EXPECT().UpdateUserCoins(gomock.Any(), fromUserID, fromUser.Coins-amount).Return(nil)
				mockRepo.EXPECT().UpdateUserCoins(gomock.Any(), toUserID, toUser.Coins+amount).Return(errors.New("update error"))
			},
			fromUserID:  fromUserID,
			toUserID:    toUserID,
			amount:      amount,
			expectedErr: errors.New("update error"),
		},
		{
			name: "error creating transaction",
			setup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), fromUserID).Return(fromUser, nil)
				mockRepo.EXPECT().GetUserByID(gomock.Any(), toUserID).Return(toUser, nil)
				mockRepo.EXPECT().UpdateUserCoins(gomock.Any(), fromUserID, fromUser.Coins-amount).Return(nil)
				mockRepo.EXPECT().UpdateUserCoins(gomock.Any(), toUserID, toUser.Coins+amount).Return(nil)
				mockRepo.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(errors.New("transaction error"))
			},
			fromUserID:  fromUserID,
			toUserID:    toUserID,
			amount:      amount,
			expectedErr: errors.New("transaction error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := service.SendCoins(context.Background(), tt.fromUserID, tt.toUserID, tt.amount)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
