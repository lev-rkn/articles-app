package controllers

import (
	authv1 "articles-service/api/auth-service/gen/proto"
	"articles-service/internal/clients/grpc/auth"
	"articles-service/internal/config"
	"articles-service/internal/lib/utils"
	"articles-service/internal/models"
	"context"
	"encoding/json"
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
	router.POST("/refresh/", authController.Refresh)

	return authController
}

//	@Summary	Аутентификация пользователя
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		user	body		models.User	true	"Почта и пароль пользователя"
//	@Success	200		{string}	string		"token"
//	@Failure	400		{string}	string		"Bad Request"
//	@Failure	500		{string}	string		"Internal Server Error"
//	@Router		/user/login/ [post]
func (a *authController) Login(c *gin.Context) {
	user := &models.User{}
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		utils.ErrorLog("unable to decode ad", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loginIn := &authv1.LoginRequest{
		Email:    user.Email,
		Password: user.Password,
		AppId:    config.Cfg.AuthGPRC.AppId,
		Fingerprint: c.GetHeader("X-fingerprint"),
	}
	res, err := a.authClient.Api.Login(c, loginIn)
	if err != nil {
		utils.ErrorLog("unable to login", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": res.GetAccessToken(),
		"refresh_token": res.GetRefreshToken(),
	})
}

//	@Summary	Регистрация пользователя
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		user	body		models.User	true	"Почта и пароль пользователя"
//	@Success	201		{string}	int			"id"
//	@Failure	400		{string}	string		"Bad Request"
//	@Failure	500		{string}	string		"Internal Server Error"
//	@Router		/user/register/ [post]
func (a *authController) Register(c *gin.Context) {
	user := &models.User{}
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		utils.ErrorLog("unable to decode user", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	registerIn := &authv1.RegisterRequest{
		Email:    user.Email,
		Password: user.Password,
	}
	res, err := a.authClient.Api.Register(c, registerIn)
	if err != nil {
		utils.ErrorLog("unable to register", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": res.GetUserId()})
}

func (a *authController) Refresh(c *gin.Context) {
	refreshRequest := struct {
		RefreshToken string `json:"refresh_token"`
	}{}
	err := json.NewDecoder(c.Request.Body).Decode(&refreshRequest)
	if err != nil {
		utils.ErrorLog("unable to decode ad", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refershIn := &authv1.RefreshTokenRequest{
		AppId:    config.Cfg.AuthGPRC.AppId,
		Fingerprint: c.GetHeader("X-fingerprint"),
		RefreshToken: refreshRequest.RefreshToken,
	}
	res, err := a.authClient.Api.Refresh(c, refershIn)
	if err != nil {
		utils.ErrorLog("unable to refresh token pair", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}