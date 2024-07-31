package controllers

import (
	"ads-service/internal/service"
	"context"
	"net/http"

	_ "ads-service/docs"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/http-swagger"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
func New(ctx context.Context, service *services.Service) *http.ServeMux {	
	mux := http.NewServeMux()

	// swagger
	mux.Handle("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	
	// Метрики
	mux.Handle("/metrics", promhttp.Handler())

	// Инициализация Контроллера объявлений
	InitAdController(ctx, service.Ad, mux)

	return mux
}
