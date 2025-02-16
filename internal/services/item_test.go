package services

import (
	"context"
	"errors"
	"testing"

	"github.com/derticom/merch-store/internal/models"
	"github.com/derticom/merch-store/internal/services/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_GetAllItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_services.NewMockRepository(ctrl)
	service := New(mockRepo)

	items := []models.Item{
		{Name: "item1", Price: 100},
		{Name: "item2", Price: 200},
	}

	tests := []struct {
		name        string
		setup       func()
		expected    []models.Item
		expectedErr error
	}{
		{
			name: "successful get all items",
			setup: func() {
				mockRepo.EXPECT().GetAllItems(gomock.Any()).Return(items, nil)
			},
			expected:    items,
			expectedErr: nil,
		},
		{
			name: "error getting all items",
			setup: func() {
				mockRepo.EXPECT().GetAllItems(gomock.Any()).Return(nil, errors.New("repository error"))
			},
			expected:    nil,
			expectedErr: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			result, err := service.GetAllItems(context.Background())
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestService_GetItemByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_services.NewMockRepository(ctrl)
	service := New(mockRepo)

	itemName := "item1"
	item := &models.Item{Name: itemName, Price: 100}

	tests := []struct {
		name        string
		setup       func()
		itemName    string
		expected    *models.Item
		expectedErr error
	}{
		{
			name: "successful get item by name",
			setup: func() {
				mockRepo.EXPECT().GetItemByName(gomock.Any(), itemName).Return(item, nil)
			},
			itemName:    itemName,
			expected:    item,
			expectedErr: nil,
		},
		{
			name: "item not found",
			setup: func() {
				mockRepo.EXPECT().GetItemByName(gomock.Any(), itemName).Return(nil, nil)
			},
			itemName:    itemName,
			expected:    nil,
			expectedErr: nil,
		},
		{
			name: "error getting item by name",
			setup: func() {
				mockRepo.EXPECT().GetItemByName(gomock.Any(), itemName).Return(nil, errors.New("repository error"))
			},
			itemName:    itemName,
			expected:    nil,
			expectedErr: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			result, err := service.GetItemByName(context.Background(), tt.itemName)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
