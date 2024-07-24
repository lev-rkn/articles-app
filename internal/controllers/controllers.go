package controllers

import (
	"ads-service/internal/service"
	"context"
	"log/slog"
	"net/http"
)

func New(service *services.Service, ctx context.Context, logger *slog.Logger) *http.ServeMux {
	mux := http.NewServeMux()
    adController := NewAdController(service, ctx, logger)
    mux.HandleFunc("POST /ad/create/", adController.Create)
    mux.HandleFunc("GET /ad/all/", adController.GetAll)
    mux.HandleFunc("GET /ad/{id}", adController.GetOne)
	
	return mux
}