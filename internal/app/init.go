package app

import (
	"ads-service/internal/config"
	"ads-service/internal/controllers"
	"ads-service/internal/service/ad"
	"ads-service/internal/store"
	"ads-service/logger"
	"context"

	"net/http"

	_ "github.com/lib/pq"
)

func Initialization() {
	// init logger
	logger := logger.New()

	// init config
	cfg := config.New(logger)

	// init db
	store := store.NewRepository(logger, cfg)

	// init router
	router := controllers.New(ad.NewService(store, logger), context.Background(), logger)

	// run server
	logger.Info("start server",
		"address", cfg.MustString("http_server_address"))
	http.ListenAndServe(cfg.MustString("http_server_address"), router)
}
