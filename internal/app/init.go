package app

import (
	"ads-service/internal/config"
	"ads-service/internal/controllers"
	"ads-service/internal/repository"
	services "ads-service/internal/service"
	"ads-service/logger"
	"context"

	"net/http"

	_ "github.com/lib/pq"
)

func Initialization() {
	// Иницилизация логгера
	logger := logger.New()
	if err := recover(); err != nil {
		logger.Error("panic", "err", err)
	}

	// Иницилизация конфига
	cfg := config.New(logger)

	// Иницилизация хранилища
	repository := repository.NewRepository(logger, cfg)

	// Новый конкекст
	ctx := context.Background()

	// Иницилизация контроллеров
	router := controllers.New(ctx, services.NewService(repository, logger), logger)

	// Запуск сервера
	serverAddr := cfg.String("http_server_address")
	logger.Info("starting server at port: " + serverAddr)
	
	if err := http.ListenAndServe(serverAddr, router); err != nil {
		logger.Error("unable to start server", "err", err.Error())
		return
	}
}
