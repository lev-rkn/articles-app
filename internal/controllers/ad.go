package controllers

import (
	"ads-service/internal/models"
	services "ads-service/internal/service"
	"ads-service/metrics"
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
	ctx       context.Context
	adService services.AdServiceInterface
	mux       *http.ServeMux
}

func InitAdController(
	ctx context.Context,
	adService services.AdServiceInterface,
	mux *http.ServeMux,
) *adController {
	adController := &adController{
		ctx:       ctx,
		adService: adService,
		mux:       mux,
	}
	mux.HandleFunc("POST /ad/create/", adController.Create)
	mux.HandleFunc("GET /ad/all/", adController.GetAll)
	mux.HandleFunc("GET /ad/{id}", adController.GetOne)

	return adController
}

var (
	ErrEmptyTitle         = errors.New("заголовок не может быть пустым")
	ErrInvalidTitle       = errors.New("длина заголовка должна превышать 200 символов")
	ErrInvalidDescription = errors.New("длина описания не должна превышать 1000 симоволов")
	ErrInvalidPhotos      = errors.New("нельзя загрузить больше чем 3 ссылки на фото")
	ErrInvalidPageNumber  = errors.New("невалидный номер страницы")
	ErrInvalidPriceSort   = errors.New("невалидный параметр сортировки по цене")
	ErrInvalidDateSort    = errors.New("невалидный параметр сортировки по дате")
	ErrInvalidId          = errors.New("невалидный идентификатор (id) объявления")
)

// @Summary Создание объявления
// @Tags ads
// @Accept json
// @Produce json
// @Param ad body models.Ad true "Объявление"
// @Success 201 {object} models.Ad
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router			/ad/create/ [post]
func (h *adController) Create(w http.ResponseWriter, r *http.Request) {
	var err error
	// собираем метрики
	go metrics.CreateAdRequest.Inc()
	defer func() {
		if err == nil {
			go metrics.CreateAdOK.Inc()
		} else {
			go metrics.CreateAdError.Inc()
		}
	}()

	ad := &models.Ad{}
	err = json.NewDecoder(r.Body).Decode(&ad)
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

	id, err := h.adService.Create(ad)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"id": %d}`, id)))
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
func (h *adController) GetAll(w http.ResponseWriter, r *http.Request) {
	var err error
	// собираем метрики
	go metrics.GetAdsRequest.Inc()
	defer func() {
		if err == nil {
			go metrics.GetAdsOK.Inc()
		} else {
			go metrics.GetAdsError.Inc()
		}
	}()

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

	adsArr, err := h.adService.GetAll(price, date, pageN)
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

// @Summary Получение одного объявления
// @Tags ads
// @Accept json
// @Produce json
// @Param id path int true "ID объявления"
// @Success 200 {object} models.Ad
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router			/ad/{id} [get]
func (h *adController) GetOne(w http.ResponseWriter, r *http.Request) {
	var err error
	// собираем метрики
	go metrics.GetAdRequest.Inc()
	defer func() {
		if err == nil {
			go metrics.GetAdOK.Inc()
		} else {
			go metrics.GetAdError.Inc()
		}
	}()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("unable to parse id", "err", err.Error())
		http.Error(w, ErrInvalidId.Error(), http.StatusBadRequest)
		return
	}

	ad, err := h.adService.GetOne(id)
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
