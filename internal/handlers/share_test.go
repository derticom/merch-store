package handlers

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestHandler_SendCoin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_handlers.NewMockService(ctrl)
	handler := New(mockService)

	fromUserID := uuid.New()
	toUserID := uuid.New()
	amount := 100

	tests := []struct {
		name           string
		setup          func()
		requestBody    map[string]interface{}
		userID         uuid.UUID
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful send coin",
			setup: func() {
				mockService.EXPECT().SendCoins(gomock.Any(), fromUserID, toUserID, amount).Return(nil)
			},
			requestBody: map[string]interface{}{
				"toUser": toUserID.String(),
				"amount": amount,
			},
			userID:         fromUserID,
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name:           "invalid request body",
			setup:          func() {},
			requestBody:    nil,
			userID:         fromUserID,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid request body\n",
		},
		{
			name:  "invalid user ID",
			setup: func() {},
			requestBody: map[string]interface{}{
				"toUser": "invalid-uuid",
				"amount": amount,
			},
			userID:         fromUserID,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid user ID\n",
		},
		{
			name: "insufficient coins",
			setup: func() {
				mockService.EXPECT().SendCoins(gomock.Any(), fromUserID, toUserID, amount).Return(errors.New("insufficient coins"))
			},
			requestBody: map[string]interface{}{
				"toUser": toUserID.String(),
				"amount": amount,
			},
			userID:         fromUserID,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "insufficient coins\n",
		},
		{
			name: "user not found",
			setup: func() {
				mockService.EXPECT().SendCoins(gomock.Any(), fromUserID, toUserID, amount).Return(errors.New("user not found"))
			},
			requestBody: map[string]interface{}{
				"toUser": toUserID.String(),
				"amount": amount,
			},
			userID:         fromUserID,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "user not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			var reqBody []byte
			if tt.requestBody != nil {
				reqBody, _ = json.Marshal(tt.requestBody)
			}

			req, err := http.NewRequest(http.MethodPost, "/send", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)

			ctx := context.WithValue(req.Context(), userIDKey, tt.userID)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/send", handler.SendCoin).Methods(http.MethodPost)

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}
