package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/derticom/merch-store/config"
	"github.com/derticom/merch-store/internal/handlers"
	"github.com/derticom/merch-store/internal/repositories"
	"github.com/derticom/merch-store/internal/server"
	"github.com/derticom/merch-store/internal/services"

	"github.com/gorilla/mux"
)

func Run(ctx context.Context, cfg *config.Config, log *slog.Logger) error {
	storage, err := repositories.New(ctx, cfg.PostgresURL)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := storage.Migrate("./migrations"); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	service := services.New(storage)

	handler := handlers.New(service)

	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	srv := server.New(cfg.Port, router, log)

	return srv.Start()
}
