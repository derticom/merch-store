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
	// Инициализация хранилища
	storage, err := repositories.New(ctx, cfg.PostgresURL)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Применение миграций
	if err := storage.Migrate("./migrations"); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Инициализация сервисов
	service := services.New(storage)

	// Инициализация обработчиков
	handler := handlers.NewHandler(service)

	// Настройка маршрутизатора
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	// Создание и запуск сервера
	srv := server.New(cfg.Port, router, log)

	return srv.Start()
}
