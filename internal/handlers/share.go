package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// SendCoin - обработчик для передачи монет.
func (h *Handler) SendCoin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	fromUserID := r.Context().Value(userIDKey).(uuid.UUID)
	toUserID, err := uuid.Parse(req.ToUser)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.SendCoins(r.Context(), fromUserID, toUserID, req.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
