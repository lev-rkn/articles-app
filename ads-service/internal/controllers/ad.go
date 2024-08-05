package controllers

import (
	"ads-service/internal/controllers/middlewares"
	"ads-service/internal/lib/types"
	"ads-service/internal/models"
	services "ads-service/internal/service"
	"ads-service/metrics"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type adController struct {
	ctx       context.Context
	adService services.AdServiceInterface
	router    *gin.RouterGroup
}

func InitAdController(
	ctx context.Context,
	adService services.AdServiceInterface,
	router *gin.RouterGroup,
) *adController {
	adController := &adController{
		ctx:       ctx,
		adService: adService,
		router:    router,
	}
	router.GET("/all/", adController.GetAll)
	router.GET("/:id", adController.GetOne)
	router.Use(middlewares.AuthMiddleware())
	router.POST("/create/", adController.Create)

	return adController
}

// @Summary Создание объявления
// @Tags ads
// @Accept json
// @Produce json
// @Param ad body models.Ad true "Объявление"
// @Success 201 {int} id
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router			/ad/create/ [post]
func (h *adController) Create(c *gin.Context) {
	go metrics.CreateAdRequest.Inc()
	// проверяем наличие ошибки, возможно переданной нам через middleware
	if v, ok := c.Value(types.KeyError).(error); ok {
		if v != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": v.Error()})
			return
		}
	}

	var err error
	defer func() {
		if err == nil {
			go metrics.CreateAdOK.Inc()
		} else {
			go metrics.CreateAdError.Inc()
		}
	}()

	ad := &models.Ad{}
	err = json.NewDecoder(c.Request.Body).Decode(&ad)
	if err != nil {
		slog.Error("unable to decode ad", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// достаем токен из контекста и берем оттуда id пользователя, создавшего это объявление
	if v, ok := c.Value(types.KeyUser).(*jwt.Token); ok {
		if claims, ok := v.Claims.(jwt.MapClaims); ok {
			if idF, ok := claims["uid"].(float64); ok {
				id := int(idF)
				ad.UserId = id
			}
		}
	}

	err = validateCreateAdRequest(ad)
	if err != nil {

	}

	id, err := h.adService.Create(ad)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// @Summary Получение страницы объявлений
// @Tags ads
// @Accept json
// @Produce json
// @Param page query int true "Номер страницы"
// @Param price query string false "Сортировка по цене (asc, desc)"
// @Param date query string false "Сортировка по дате (asc, desc)"
// @Success 200 {array} models.Ad
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router			/ad/all/ [get]
func (h *adController) GetAll(c *gin.Context) {
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

	q := c.Request.URL.Query()
	page, price, date, userId := c.Query("page"), c.Query("price"), c.Query("date"), c.Query("user_id")

	// проверка, что номер страницы является целочисленным значением
	pageN, err := strconv.Atoi(page)
	if err != nil {
		slog.Error("unable to parse page number", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": types.ErrInvalidPageNumber.Error()})
		return
	}

	// проверка, что идентификатор пользователя является целочисленным значением
	var userIdN int
	if q.Has("user_id") {
		// проверка, что идентификатор пользователя является целочисленным значением
		userIdN, err = strconv.Atoi(userId)
		if err != nil {
			slog.Error("unable to parse user id", "err", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": types.ErrInvalidUserId.Error()})
			return
		}
	}

	// Проверка, что параметры сортировки по цене содержат либо asc, либо desc
	if q.Has("price") {
		if price == "asc" || price == "desc" {
			price = "price " + price
		} else {
			slog.Error("Invalid price query parameter: " + price)
			c.JSON(http.StatusBadRequest, types.ErrInvalidPriceSort.Error())
			return
		}
	}

	// Проверка, что параметры сортировки по дате содержат либо asc, либо desc
	if q.Has("date") {
		if date == "asc" || date == "desc" {
			date = "date " + date
		} else {
			slog.Error("Invalid date query parameter: " + date)
			c.JSON(http.StatusBadRequest, types.ErrInvalidDateSort.Error())
			return
		}
	}

	adsArr, err := h.adService.GetAll(price, date, pageN, userIdN)
	if err != nil {
		slog.Error("unable to get ads", "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, adsArr)
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
func (h *adController) GetOne(c *gin.Context) {
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

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Error("unable to parse id", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": types.ErrInvalidId.Error()})
		return
	}

	ad, err := h.adService.GetOne(id)
	if err != nil {
		slog.Error("unable to get ad", "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		slog.Error("ad marshalling", "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ad)
}

func validateCreateAdRequest(ad *models.Ad) error {
	// проверка, что длина заголовка не превышает 200 симоволов
	if utf8.RuneCountInString(ad.Title) > 200 {
		slog.Error("invalid title", "title", ad.Title)
		return types.ErrInvalidTitle
	}
	// Проверка, что заголовок не пустой
	if ad.Title == "" {
		slog.Error("empty title", "title", ad.Title)
		return types.ErrEmptyTitle
	}
	// проверка, что длина описания не должна превышать 1000 символов
	if utf8.RuneCountInString(ad.Description) > 1000 {
		slog.Error("invalid description", "description", ad.Description)
		return types.ErrInvalidDescription
	}
	// проверка, что нельзя загрузить больше чем 3 ссылки на фото
	if len(ad.Photos) > 3 {
		slog.Error("invalid photos", "photos", ad.Photos)
		return types.ErrInvalidPhotos
	}
	return nil
}
