package handlers

import "github.com/derticom/merch-store/internal/services"

type Handler struct {
	Service services.Service
}

func NewHandler(service services.Service) *Handler {
	return &Handler{Service: service}
}
