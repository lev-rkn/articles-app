package app

import (
	"ads-service/internal/config"
	"ads-service/internal/handlers"
	"ads-service/internal/store"
	"ads-service/logger"

	"net/http"

	_ "github.com/lib/pq"
)

func Initialization() {
	// init logger
	logger := logger.New()

	// init config
	cfg := config.New(logger)

	// init db
	store.New(logger, cfg)

	// init router
	router := handlers.New()

	// run server
	logger.Info("start server",
		"address", cfg.MustString("http_server_address"))
	http.ListenAndServe(cfg.MustString("http_server_address"), router)
}
