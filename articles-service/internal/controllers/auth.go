package controllers

import (
	authv1 "articles-service/api/auth-service/gen/proto"
	"articles-service/internal/clients/grpc/auth"
	"articles-service/internal/config"
	"articles-service/internal/models"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authController struct {
	ctx        context.Context
	router     *gin.RouterGroup
	authClient *auth.Client
}

func InitAuthController(
	ctx context.Context,
	router *gin.RouterGroup,
	authClient *auth.Client,
) *authController {
	authController := &authController{
		ctx:        ctx,
		router:     router,
		authClient: authClient,
	}
	router.POST("/login/", authController.Login)
	router.POST("/register/", authController.Register)

	return authController
}

// @Summary Аутентификация пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "Почта и пароль пользователя"
// @Success 200 {string} string "token"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router	/user/login/ [post]
func (a *authController) Login(c *gin.Context) {
	user := &models.User{}
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		slog.Error("unable to decode ad", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loginIn := &authv1.LoginRequest{
		AppId:    config.Cfg.AuthGPRC.AppId,
		Email:    user.Email,
		Password: user.Password,
	}
	res, err := a.authClient.Api.Login(c, loginIn)
	if err != nil {
		slog.Error("unable to login", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": res.GetToken()})
}

// @Summary Регистрация пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "Почта и пароль пользователя"
// @Success 201 {string} int "id"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router	/user/register/ [post]
func (a *authController) Register(c *gin.Context) {
	user := &models.User{}
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		slog.Error("unable to decode ad", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	registerIn := &authv1.RegisterRequest{
		Email:    user.Email,
		Password: user.Password,
	}
	res, err := a.authClient.Api.Register(c, registerIn)
	if err != nil {
		slog.Error("unable to register", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": res.GetUserId()})
}
