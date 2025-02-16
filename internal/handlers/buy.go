package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// BuyItem - обработчик для покупки товара.
func (h *Handler) BuyItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemName := vars["item"]

	userID := r.Context().Value(userIDKey).(uuid.UUID)

	if err := h.service.BuyItem(r.Context(), userID, itemName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
