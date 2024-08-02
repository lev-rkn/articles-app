package controllers

import (
	authv1 "ads-service/api/auth-service/gen/proto"
	"ads-service/internal/clients/grpc/auth"
	"ads-service/internal/config"
	"ads-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type userController struct {
	ctx        context.Context
	mux        *http.ServeMux
	authClient *auth.Client
}

func InitUserController(
	ctx context.Context,
	mux *http.ServeMux,
	authClient *auth.Client,
) *userController {
	userController := &userController{
		ctx:        ctx,
		mux:        mux,
		authClient: authClient,
	}
	mux.HandleFunc("POST /user/create/", userController.Login)
	mux.HandleFunc("POST /user/register/", userController.Register)

	return userController
}

// @Summary Создание объявления
// @Tags ads
// @Accept json
// @Produce json
// @Param ad body models.Ad true "Объявление"
// @Success 201 {object} models.Ad
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router			/ad/create/ [post]
func (c *userController) Login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		slog.Error("unable to decode ad", "err", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loginIn := &authv1.LoginRequest{
		AppId:    config.Cfg.AuthGPRC.AppId,
		Email:    user.Email,
		Password: user.Password,
	}
	res, err := c.authClient.Api.Login(c.ctx, loginIn)
	if err != nil {
		slog.Error("unable to login", "err", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, res.GetToken())))
}

// @Summary Получение страницы объявлений
// @Tags ads
// @Accept json
// @Produce json
// @Param page query int true "Номер страницы"
// @Param price query string false "Сортировка по цене (asc, desc)"
// @Param date query string false "Сортировка по дате (asc, desc)"
// @Success 200 {object} models.Ad
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router			/ad/all/ [get]
func (c *userController) Register(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		slog.Error("unable to decode ad", "err", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	registerIn := &authv1.RegisterRequest{
		Email:    user.Email,
		Password: user.Password,
	}
	res, err := c.authClient.Api.Register(c.ctx, registerIn)
	if err != nil {
		slog.Error("unable to register", "err", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"id": %d}`, res.GetUserId())))
}
