package controllers

import (
	"articles-service/internal/controllers/middlewares"
	"articles-service/internal/lib/types"
	"articles-service/internal/lib/utils"
	"articles-service/internal/models"
	"articles-service/internal/service"
	"articles-service/metrics"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
)

type articleController struct {
	ctx            context.Context
	articleService service.ArticleServiceInterface
	router         *gin.RouterGroup
}

func InitArticleController(
	ctx context.Context,
	articleService service.ArticleServiceInterface,
	router *gin.RouterGroup,
) *articleController {
	articleController := &articleController{
		ctx:            ctx,
		articleService: articleService,
		router:         router,
	}
	router.GET("/all/", articleController.GetAll)
	router.GET("/:id", articleController.GetOne)
	router.Use(middlewares.AuthMiddleware())
	router.POST("/create/", articleController.Create)

	return articleController
}

// @Summary Создание объявления
// @Tags articles
// @Accept json
// @Produce json
// @Param article body models.Article true "Объявление"
// @Success 201 {int} id
// @Failure 400 {string} string "Barticle Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router			/article/create/ [post]
func (h *articleController) Create(c *gin.Context) {
	go metrics.CreateArticleRequest.Inc()
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
			go metrics.CreateArticleOK.Inc()
		} else {
			go metrics.CreateArticleError.Inc()
		}
	}()

	article := &models.Article{}
	err = json.NewDecoder(c.Request.Body).Decode(&article)
	if err != nil {
		utils.ErrorLog("unable to decode article", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// достаем токен из контекста и берем оттуда id пользователя, создавшего это объявление
	if v, ok := c.Value(types.KeyUser).(*jwt.Token); ok {
		if claims, ok := v.Claims.(jwt.MapClaims); ok {
			if idF, ok := claims["uid"].(float64); ok {
				id := int(idF)
				article.UserId = id
			}
		}
	}

	validate := validator.New()
	err = validate.Struct(article)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.articleService.Create(article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// @Summary Получение страницы объявлений
// @Tags articles
// @Accept json
// @Produce json
// @Param page query int true "Номер страницы"
// @Param price query string false "Сортировка по цене (asc, desc)"
// @Param date query string false "Сортировка по дате (asc, desc)"
// @Success 200 {array} models.Article
// @Failure 400 {string} string "Barticle Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router			/article/all/ [get]
func (h *articleController) GetAll(c *gin.Context) {
	var err error
	// собираем метрики
	go metrics.GetArticlesRequest.Inc()
	defer func() {
		if err == nil {
			go metrics.GetArticlesOK.Inc()
		} else {
			go metrics.GetArticlesError.Inc()
		}
	}()

	page,
		price,
		date,
		userId := c.Query("page"),
		c.Query("price"),
		c.Query("date"),
		c.Query("user_id")

	// проверка, что номер страницы является целочисленным значением
	pageN, err := strconv.Atoi(page)
	if err != nil {
		utils.ErrorLog("unable to parse page number", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": types.ErrInvalidPageNumber.Error()})
		return
	}

	// проверка, что идентификатор пользователя является целочисленным значением
	var userIdN int
	if userId != "" {
		// проверка, что идентификатор пользователя является целочисленным значением
		userIdN, err = strconv.Atoi(userId)
		if err != nil {
			utils.ErrorLog("unable to parse user id", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": types.ErrInvalidUserId.Error()})
			return
		}
	}

	// Проверка, что параметры сортировки по цене содержат либо asc, либо desc
	if price != "" {
		if price == "asc" || price == "desc" {
			price = "price " + price
		} else {
			utils.ErrorLog("Invalid price query parameter: "+price, types.ErrInvalidPriceSort)
			c.JSON(http.StatusBadRequest, types.ErrInvalidPriceSort.Error())
			return
		}
	}

	// Проверка, что параметры сортировки по дате содержат либо asc, либо desc
	if date != "" {
		if date == "asc" || date == "desc" {
			date = "date " + date
		} else {
			utils.ErrorLog("Invalid date query parameter: "+date, types.ErrInvalidDateSort)
			c.JSON(http.StatusBadRequest, types.ErrInvalidDateSort.Error())
			return
		}
	}

	articlesArr, err := h.articleService.GetAll(price, date, pageN, userIdN)
	if err != nil {
		utils.ErrorLog("unable to get articles", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, articlesArr)
}

// @Summary Получение одного объявления
// @Tags articles
// @Accept json
// @Produce json
// @Param id path int true "ID объявления"
// @Success 200 {object} models.Article
// @Failure 400 {string} string "Barticle Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router			/article/{id} [get]
func (h *articleController) GetOne(c *gin.Context) {
	var err error
	// собираем метрики
	go metrics.GetArticleRequest.Inc()
	defer func() {
		if err == nil {
			go metrics.GetArticleOK.Inc()
		} else {
			go metrics.GetArticleError.Inc()
		}
	}()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorLog("unable to parse id", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": types.ErrInvalidId.Error()})
		return
	}

	article, err := h.articleService.GetOne(id)
	if err != nil {
		utils.ErrorLog("unable to get article", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		utils.ErrorLog("article marshalling", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, article)
}