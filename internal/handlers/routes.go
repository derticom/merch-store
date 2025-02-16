package handlers

import (
	"github.com/gorilla/mux"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/auth/register", h.Register).Methods("POST")
	api.HandleFunc("/auth/login", h.Login).Methods("POST")

	api.HandleFunc("/info", h.GetInfo).Methods("GET")
	api.HandleFunc("/sendCoin", h.SendCoin).Methods("POST")
	api.HandleFunc("/buy/{item}", h.BuyItem).Methods("GET")
}
