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
	"unicode/utf8"
)

type adController struct {
	service *services.Service
	ctx       context.Context
	logger    *slog.Logger
}

func NewAdController(
	ctx context.Context,
	service *services.Service,
	logger *slog.Logger,
) *adController {
	return &adController{
		ctx:       ctx,
		service: service,
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

	// проверка, что длина заголовка от 1 до 200
	if len(ad.Title) > 0 || utf8.RuneCountInString(ad.Title) > 200 {
		h.logger.Error("invalid title", "title", ad.Title)
		http.Error(w, "Длина заголовка должна быть длиной от 1 до 200 симоволов", http.StatusBadRequest)
		return
	}
	// проверка, что длина описания не должна превышать 1000 символов
	if utf8.RuneCountInString(ad.Description) > 1000 {
		h.logger.Error("invalid description", "description", ad.Description)
		http.Error(w, "Длина описания должна быть длиной до 1000 симоволов", http.StatusBadRequest)
		return
	}
	// проверка, что нельзя загрузить больше чем 3 ссылки на фото
	if len(ad.Photos) > 3 {
		h.logger.Error("invalid photos", "photos", ad.Photos)
		http.Error(w, "Нельзя загрузить больше чем 3 ссылки на фото", http.StatusBadRequest)
		return
	}

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

	// проверка, что номер страницы является целочисленным значением
	pageN, err := strconv.Atoi(page)
	if err != nil {
		h.logger.Error("unable to parse page number", "err", err.Error())
		http.Error(w, "Невалидный номер страницы", http.StatusBadRequest)
		return
	}

	// Проверка, что параметры сортировки по цене содержат либо asc, либо desc
	if q.Has("price") {
		if price == "asc" || price == "desc" {
			price = "price " + price
		} else {
			h.logger.Error("Invalid price query parameter: " + price)
			http.Error(w, "Invalid price query parameter: "+price, http.StatusBadRequest)
			return
		}
	}

	// Проверка, что параметры сортировки по дате содержат либо asc, либо desc
	if q.Has("date") {
		if date == "asc" || date == "desc" {
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
