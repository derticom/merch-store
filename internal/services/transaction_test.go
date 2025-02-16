package services

import (
	"context"
	"errors"
	"testing"

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
				mockRepo.EXPECT().SendCoins(gomock.Any(), fromUserID, toUserID, amount).Return(nil)
			},
			fromUserID:  fromUserID,
			toUserID:    toUserID,
			amount:      amount,
			expectedErr: nil,
		},
		{
			name:        "sender and receiver are the same",
			setup:       func() {},
			fromUserID:  fromUserID,
			toUserID:    fromUserID,
			amount:      amount,
			expectedErr: errors.New("cannot send coins to yourself"),
		},
		{
			name:        "non-positive amount",
			setup:       func() {},
			fromUserID:  fromUserID,
			toUserID:    toUserID,
			amount:      0,
			expectedErr: errors.New("amount must be positive"),
		},
		{
			name: "error in repository",
			setup: func() {
				mockRepo.EXPECT().SendCoins(gomock.Any(), fromUserID, toUserID, amount).Return(errors.New("repository error"))
			},
			fromUserID:  fromUserID,
			toUserID:    toUserID,
			amount:      amount,
			expectedErr: errors.New("repository error"),
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
