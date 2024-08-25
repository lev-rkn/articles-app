package controllers

import (
	"articles-service/internal/clients/grpc/auth"
	"articles-service/internal/config"
	"articles-service/internal/controllers/middlewares"
	"articles-service/internal/lib/utils"
	"articles-service/internal/service"
	"context"
	"time"

	_ "articles-service/docs"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
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

// @host	localhost:8080
func NewRouter(ctx context.Context, service *service.Service) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.AuthMiddleware())
	// TODO: обновить документацию эндпоинтов swagger
	router.GET("/swagger/*any", gin.WrapF(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	)))

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	articleRouter := router.Group("/article/")
	InitArticleController(ctx, service.Article, articleRouter)

	commentRouter := router.Group("/comments/")
	InitCommentController(ctx, service.Comment, commentRouter)

	// создание экземпляра клиента для обращения к auth по GRPC
	authClient, err := auth.NewAuthClient(ctx, config.Cfg.AuthGPRC.Address, time.Second*10, 3)
	if err != nil {
		utils.ErrorLog("create authGRPC client", err)
	}
	authRouter := router.Group("/auth")
	InitAuthController(ctx, authRouter, authClient)

	return router
}
