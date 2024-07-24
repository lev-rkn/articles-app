package app

import (
	"ads-service/internal/config"
	"ads-service/internal/controllers"
	services "ads-service/internal/service"
	"ads-service/internal/store"
	"ads-service/logger"
	"context"

	"net/http"

	_ "github.com/lib/pq"
)

func Initialization() {
	// Иницилизация логгера
	logger := logger.New()

	// Иницилизация конфига
	cfg := config.New(logger)

	// Иницилизация хранилища
	store := store.NewRepository(logger, cfg)

	// Иницилизация контроллеров
	router := controllers.New(services.NewService(store, logger), context.Background(), logger)

	// Запуск сервера
	logger.Info("start server",
		"address", cfg.MustString("http_server_address"))
	http.ListenAndServe(cfg.MustString("http_server_address"), router)
}
