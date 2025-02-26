package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// GetInfo - обработчик для получения информации о пользователе.
func (h *Handler) GetInfo(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uuid.UUID)

	user, err := h.service.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	transactions, err := h.service.GetTransactionHistory(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get transaction history", http.StatusInternalServerError)
		return
	}

	purchases, err := h.service.GetPurchaseHistory(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get purchase history", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"coins":     user.Coins,
		"inventory": purchases,
		"coinHistory": map[string]interface{}{
			"received": transactions,
			"sent":     transactions,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
