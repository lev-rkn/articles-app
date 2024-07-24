package controllers

import (
	"ads-service/internal/models"
	"ads-service/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

type adController struct {
	service *services.Service
	ctx       context.Context
	logger    *slog.Logger
}

func NewAdController(
	service *services.Service,
	ctx context.Context,
	logger *slog.Logger,
) *adController {
	return &adController{
		service: service,
		ctx:       ctx,
		logger:    logger,
	}
}

func (h *adController) Create(w http.ResponseWriter, r *http.Request) {
	ad := &models.Ad{}
	err := json.NewDecoder(r.Body).Decode(&ad)
	if err != nil {
		h.logger.Error("unable to decode ad", "err", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: validate req

	id, err := h.service.Ad.Create(ad)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"id": %d}`, id)))
}

func (h *adController) GetAll(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, price, date := q.Get("page"), q.Get("price"), q.Get("date")

	// checking that page query parameter is an integer
	pageN, err := strconv.Atoi(page)
	if err != nil {
		h.logger.Error("unable to parse page number", "err", err.Error())
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	// checking that price query parameter equal to asc or desc
	if q.Has("price") {
		if price == "asc" || price == "desc" {
			// for build query
			price = "price " + price
		} else {
			h.logger.Error("Invalid price query parameter: " + price)
			http.Error(w, "Invalid price query parameter: "+price, http.StatusBadRequest)
			return
		}
	}

	// checking that date query parameter equal to asc or desc
	if q.Has("date") {
		if date == "asc" || date == "desc" {
			// for build query
			date = "date " + date
		} else {
			h.logger.Error("Invalid date query parameter: " + date)
			http.Error(w, "Invalid date query parameter: "+date, http.StatusBadRequest)
			return
		}
	}

	adsArr, err := h.service.Ad.GetAll(price, date, pageN)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	marshalled, err := json.Marshal(adsArr)
	if err != nil {
		h.logger.Error("adsArr marshalling", "err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshalled)
}

func (h *adController) GetOne(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.logger.Error("unable to parse id", "err", err.Error())
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	ad, err := h.service.Ad.GetOne(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	marshalled, err := json.Marshal(ad)
	if err != nil {
		h.logger.Error("ad marshalling", "err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshalled)
}
