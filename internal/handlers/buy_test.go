package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/derticom/merch-store/internal/handlers/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandler_BuyItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_handlers.NewMockService(ctrl)
	handler := New(mockService)

	userID := uuid.New()
	itemName := "umbrella"

	tests := []struct {
		name           string
		setup          func()
		itemName       string
		userID         uuid.UUID
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful purchase",
			setup: func() {
				mockService.EXPECT().BuyItem(gomock.Any(), userID, itemName).Return(nil)
			},
			itemName:       itemName,
			userID:         userID,
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name: "insufficient coins",
			setup: func() {
				mockService.EXPECT().BuyItem(gomock.Any(), userID, itemName).Return(errors.New("insufficient coins"))
			},
			itemName:       itemName,
			userID:         userID,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "insufficient coins\n",
		},
		{
			name: "item not found",
			setup: func() {
				mockService.EXPECT().BuyItem(gomock.Any(), userID, itemName).Return(errors.New("item not found"))
			},
			itemName:       itemName,
			userID:         userID,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "item not found\n",
		},
		{
			name: "user not found",
			setup: func() {
				mockService.EXPECT().BuyItem(gomock.Any(), userID, itemName).Return(errors.New("user not found"))
			},
			itemName:       itemName,
			userID:         userID,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "user not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			req, err := http.NewRequest(http.MethodPost, "/buy/"+tt.itemName, nil)
			assert.NoError(t, err)

			ctx := context.WithValue(req.Context(), userIDKey, tt.userID)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/buy/{item}", handler.BuyItem).Methods(http.MethodPost)

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}
