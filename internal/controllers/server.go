package controllers

import (
	"ads-service/internal/service"
	"context"
	"net/http"
)

func New(ctx context.Context, service *services.Service) *http.ServeMux {
	mux := http.NewServeMux()
    InitAdController(ctx, service.Ad, mux)
	
	return mux
}