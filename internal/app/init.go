package app

import (
	"ads-service/internal/config"
	"ads-service/internal/controllers"
	"ads-service/internal/repository"
	"ads-service/internal/service"
	"ads-service/logger"
	"context"
	"log/slog"

	"net/http"

	_ "github.com/lib/pq"
)

func Initialization() {
	// Иницилизация логгера
	logger.Init()

	// Иницилизация конфига
	cfg := config.New()

	// Новый конкекст
	ctx := context.Background()

	// Иницилизация хранилища
	repository := repository.NewRepository(ctx, cfg)

	// Иницилизация контроллеров
	router := controllers.New(ctx, services.NewService(repository))

	// Запуск сервера
	serverAddr := cfg.String("http_server_address")
	slog.Info("starting server at port: " + serverAddr)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		slog.Error("unable to start server", "err", err.Error())
		return
	}
}
