package app

import (
	"articles-service/internal/config"
	"articles-service/internal/controllers"
	"articles-service/internal/repository"
	"articles-service/internal/service"
	"articles-service/logger"
	"context"

	"net/http"

	_ "github.com/lib/pq"
)

func NewServer() *http.Server {
	config.MustLoad()
	logger.MustLoad()

	mainCtx := context.Background()

	repository := repository.NewRepository(mainCtx)
	service := service.NewService(repository)
	router := controllers.NewRouter(mainCtx, service)

	server := &http.Server{
		Addr:    config.Cfg.HTTPServerAddress,
		Handler: router.Handler(),
	}

	return server
}
