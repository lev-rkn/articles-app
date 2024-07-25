package controllers

import (
	"ads-service/internal/models"
	services "ads-service/internal/service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"unicode/utf8"
)

type adController struct {
	service *services.Service
	ctx     context.Context
}

func NewAdController(
	ctx context.Context,
	service *services.Service,
) *adController {
	return &adController{
		ctx:     ctx,
		service: service,
	}
}

var (
	ErrEmptyTitle = errors.New("заголовок не может быть пустым")
	ErrInvalidTitle = errors.New("длина заголовка должна превышать 200 символов")
	ErrInvalidDescription = errors.New("длина описания не должна превышать 1000 симоволов")
	ErrInvalidPhotos = errors.New("нельзя загрузить больше чем 3 ссылки на фото")
	ErrInvalidPageNumber = errors.New("невалидный номер страницы")
	ErrInvalidPriceSort = errors.New("невалидный параметр сортировки по цене")
	ErrInvalidDateSort = errors.New("невалидный параметр сортировки по дате")
	ErrInvalidId = errors.New("невалидный идентификатор (id) объявления")
)

func (h *adController) Create(w http.ResponseWriter, r *http.Request) {
	ad := &models.Ad{}
	err := json.NewDecoder(r.Body).Decode(&ad)
	if err != nil {
		slog.Error("unable to decode ad", "err", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// проверка, что длина заголовка не превышает 200 симоволов
	if utf8.RuneCountInString(ad.Title) > 200 {
		slog.Error("invalid title", "title", ad.Title)
		http.Error(w, ErrInvalidTitle.Error(), http.StatusBadRequest)
		return
	}
	// Проверка, что заголовок не пустой
	if ad.Title == "" {
		slog.Error("empty title", "title", ad.Title)
		http.Error(w, ErrEmptyTitle.Error(), http.StatusBadRequest)
		return
	}
	// проверка, что длина описания не должна превышать 1000 символов
	if utf8.RuneCountInString(ad.Description) > 1000 {
		slog.Error("invalid description", "description", ad.Description)
		http.Error(w, ErrInvalidDescription.Error(), http.StatusBadRequest)
		return
	}
	// проверка, что нельзя загрузить больше чем 3 ссылки на фото
	if len(ad.Photos) > 3 {
		slog.Error("invalid photos", "photos", ad.Photos)
		http.Error(w, ErrInvalidPhotos.Error(), http.StatusBadRequest)
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
		slog.Error("unable to parse page number", "err", err.Error())
		http.Error(w, ErrInvalidPageNumber.Error(), http.StatusBadRequest)
		return
	}

	// Проверка, что параметры сортировки по цене содержат либо asc, либо desc
	if q.Has("price") {
		if price == "asc" || price == "desc" {
			price = "price " + price
		} else {
			slog.Error("Invalid price query parameter: " + price)
			http.Error(w, ErrInvalidPriceSort.Error(), http.StatusBadRequest)
			return
		}
	}

	// Проверка, что параметры сортировки по дате содержат либо asc, либо desc
	if q.Has("date") {
		if date == "asc" || date == "desc" {
			date = "date " + date
		} else {
			slog.Error("Invalid date query parameter: " + date)
			http.Error(w, ErrInvalidDateSort.Error(), http.StatusBadRequest)
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
		slog.Error("adsArr marshalling", "err", err.Error())
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
		slog.Error("unable to parse id", "err", err.Error())
		http.Error(w, ErrInvalidId.Error(), http.StatusBadRequest)
		return
	}

	ad, err := h.service.Ad.GetOne(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	marshalled, err := json.Marshal(ad)
	if err != nil {
		slog.Error("ad marshalling", "err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshalled)
}
