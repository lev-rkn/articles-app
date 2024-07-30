package app

import (
	"ads-service/internal/config"
	"ads-service/internal/controllers"
	"ads-service/internal/logger"
	"ads-service/internal/repository"
	"ads-service/internal/service"
	"context"
	"log/slog"

	"net/http"

	_ "github.com/lib/pq"
)

func Initialization() {
	// Иницилизация конфига
	cfg := config.MustLoad()
	
	// Иницилизация логгера
	logger.MustLoad(cfg.CfgType)

	// Новый конкекст
	ctx := context.Background()

	// Иницилизация хранилища
	repository := repository.NewRepository(ctx, cfg)

	// Иницилизация контроллеров
	router := controllers.New(ctx, services.NewService(repository))

	// Запуск сервера
	serverAddr := cfg.HTTPServerAddress
	slog.Info("starting server at port: " + serverAddr)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		slog.Error("unable to start server", "err", err.Error())
		return
	}
}
