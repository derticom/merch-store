package handlers

import (
	"github.com/gorilla/mux"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/auth/register", h.Register).Methods("POST")
	api.HandleFunc("/auth/login", h.Login).Methods("POST")

	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(JWTAuthMiddleware)
	protected.HandleFunc("/info", h.GetInfo).Methods("GET")
	protected.HandleFunc("/sendCoin", h.SendCoin).Methods("POST")
	protected.HandleFunc("/buy/{item}", h.BuyItem).Methods("GET")
}
