package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/derticom/merch-store/internal/handlers/mocks"
	"github.com/derticom/merch-store/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_handlers.NewMockService(ctrl)
	handler := New(mockService)

	userID := uuid.New()
	user := &models.User{
		ID:    userID,
		Coins: 1000,
	}
	transactions := []models.Transaction{
		{ID: uuid.New(), FromUser: userID, ToUser: uuid.New(), Amount: 100},
	}
	purchases := []models.Purchase{
		{ID: uuid.New(), UserID: userID, Item: "item"},
	}

	tests := []struct {
		name           string
		setup          func()
		userID         uuid.UUID
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful get info",
			setup: func() {
				mockService.EXPECT().GetUserByID(gomock.Any(), userID).Return(user, nil)
				mockService.EXPECT().GetTransactionHistory(gomock.Any(), userID).Return(transactions, nil)
				mockService.EXPECT().GetPurchaseHistory(gomock.Any(), userID).Return(purchases, nil)
			},
			userID:         userID,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"coins":     user.Coins,
				"inventory": purchases,
				"coinHistory": map[string]interface{}{
					"received": transactions,
					"sent":     transactions,
				},
			},
		},
		{
			name: "user not found",
			setup: func() {
				mockService.EXPECT().GetUserByID(gomock.Any(), userID).Return(
					nil, errors.New("user not found"))
			},
			userID:         userID,
			expectedStatus: http.StatusNotFound,
			expectedBody:   nil,
		},
		{
			name: "failed to get transaction history",
			setup: func() {
				mockService.EXPECT().GetUserByID(gomock.Any(), userID).Return(user, nil)
				mockService.EXPECT().GetTransactionHistory(gomock.Any(), userID).Return(
					nil, errors.New("failed to get transaction history"))
			},
			userID:         userID,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
		{
			name: "failed to get purchase history",
			setup: func() {
				mockService.EXPECT().GetUserByID(gomock.Any(), userID).Return(user, nil)
				mockService.EXPECT().GetTransactionHistory(gomock.Any(), userID).Return(transactions, nil)
				mockService.EXPECT().GetPurchaseHistory(gomock.Any(), userID).Return(
					nil, errors.New("failed to get purchase history"))
			},
			userID:         userID,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			req, err := http.NewRequest(http.MethodGet, "/info", nil)
			assert.NoError(t, err)

			ctx := context.WithValue(req.Context(), userIDKey, tt.userID)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/info", handler.GetInfo).Methods(http.MethodGet)

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				expectedBytes, _ := json.Marshal(tt.expectedBody)
				actualBytes := rr.Body.Bytes()
				actualBytes = []byte(strings.TrimSpace(string(actualBytes)))
				assert.Equal(t, expectedBytes, actualBytes)
			}
		})
	}
}
