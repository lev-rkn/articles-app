package controllers

import (
	"ads-service/internal/clients/grpc/auth"
	"ads-service/internal/config"
	"ads-service/internal/service"
	"context"
	"log/slog"
	"time"

	_ "ads-service/docs"

	"github.com/gin-gonic/gin"
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
func New(ctx context.Context, service *services.Service) *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", gin.WrapF(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	)))

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Инициализация Контроллера объявлений
	adRouter := router.Group("/ad")
	InitAdController(ctx, service.Ad, adRouter)

	// создание экземпляра клиента для обращения к auth по GRPC
	client, err := auth.New(ctx, config.Cfg.AuthGPRC.Address, time.Second*10, 3)
	if err != nil {
		slog.Error(err.Error())
	}
	// Инициализация контроллера аутентификации-авторизации
	userRouter := router.Group("/user")
	InitAuthController(ctx, userRouter, client)

	return router
}
