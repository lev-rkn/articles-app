package controllers

import (
	"ads-service/internal/service"
	"context"
	"net/http"
)

func New(ctx context.Context, service *services.Service) *http.ServeMux {
	mux := http.NewServeMux()
    adController := NewAdController(ctx, service)
    mux.HandleFunc("POST /ad/create/", adController.Create)
    mux.HandleFunc("GET /ad/all/", adController.GetAll)
    mux.HandleFunc("GET /ad/{id}", adController.GetOne)
	
	return mux
}