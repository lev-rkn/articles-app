package app

import (
	"ads-service/internal/config"
	"ads-service/internal/controllers"
	"ads-service/internal/repository"
	"ads-service/internal/service"
	"ads-service/logger"
	"context"

	"net/http"

	_ "github.com/lib/pq"
)

func NewServer() *http.Server {
	config.MustLoad()
	logger.MustLoad()

	mainCtx := context.Background()

	repository := repository.NewRepository(mainCtx)
	service := services.NewService(repository)
	router := controllers.NewRouter(mainCtx, service)

	server := &http.Server{
		Addr:    config.Cfg.HTTPServerAddress,
		Handler: router.Handler(),
	}

	return server
}

