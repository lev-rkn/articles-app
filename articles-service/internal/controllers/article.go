package controllers

import (
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
)

type articleController struct {
	ctx            context.Context
	articleService service.ArticleServiceInterface
}

func InitArticleController(
	ctx context.Context,
	articleService service.ArticleServiceInterface,
	router *gin.RouterGroup,
) *articleController {
	articleController := &articleController{
		ctx:            ctx,
		articleService: articleService,
	}
	router.POST("/create/", articleController.CreateArticle)
	router.GET("/all/", articleController.GetAllArticles)
	router.GET("/:id", articleController.GetOneArticle)

	return articleController
}

// @Summary	Создание объявления
// @Tags		articles
// @Accept		json
// @Produce	json
// @Param		article	body		models.Article	true	"Объявление"
// @Success	201		{int}		id
// @Failure	400		{string}	string	"Barticle Request"
// @Failure	500		{string}	string	"Internal Server Error"
// @Router		/article/create/ [post]
func (h *articleController) CreateArticle(c *gin.Context) {
	go metrics.CreateArticleRequest.Inc()
	// проверяем наличие ошибки, возможно переданной нам через middleware
	if err, ok := c.Value(types.KeyError).(error); ok && err != nil {
		utils.ErrorLog("get error from context", err)
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	article := &models.Article{}
	err := json.NewDecoder(c.Request.Body).Decode(&article)
	if err != nil {
		utils.ErrorLog("unable to decode article", err)
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	parsedClaims := utils.ParseUserClaims(c.Value(types.KeyUser))
	article.UserId = parsedClaims.UserId

	validate := validator.New()
	err = validate.Struct(article)
	if err != nil {
		utils.ErrorLog("validate article", err)
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.articleService.CreateArticle(article)
	if err != nil {
		utils.ErrorLog("article creating by service", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// @Summary	Получение страницы Статей
// @Tags		articles
// @Accept		json
// @Produce	json
// @Param		page	query		int		true	"Номер страницы"
// @Param		date	query		string	false	"Сортировка по дате (asc, desc)"
// @Success	200		{array}		models.Article
// @Failure	400		{string}	string	"Barticle Request"
// @Failure	500		{string}	string	"Internal Server Error"
// @Router		/article/all/ [get]
func (h *articleController) GetAllArticles(c *gin.Context) {
	go metrics.GetArticlesRequest.Inc()

	page,
		dateSorting,
		userId := c.Query("page"),
		c.Query("date"),
		c.Query("user_id")

	// проверка, что номер страницы является целочисленным значением
	pageN, err := strconv.Atoi(page)
	if err != nil {
		utils.ErrorLog("unable to parse page number", err)
		utils.ErrorResponse(c, http.StatusBadRequest, types.ErrInvalidPageNumber)
		return
	}

	// проверка, что идентификатор пользователя является целочисленным значением
	var userIdN int
	if userId != "" {
		userIdN, err = strconv.Atoi(userId)
		if err != nil {
			utils.ErrorLog("unable to parse user id", err)
			utils.ErrorResponse(c, http.StatusBadRequest, types.ErrInvalidUserId)
			return
		}
	}

	// Проверка, что параметры сортировки по дате содержат либо asc, либо desc
	if dateSorting != "" {
		if dateSorting != "asc" && dateSorting != "desc" {
			utils.ErrorLog("Invalid date query parameter"+dateSorting, nil)
			utils.ErrorResponse(c, http.StatusBadRequest, types.ErrInvalidDateSort)
			return
		}
	}

	articlesArr, err := h.articleService.GetAllArticles(dateSorting, pageN, userIdN)
	if err != nil {
		utils.ErrorLog("unable to get articles", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, articlesArr)
}

// @Summary	Получение одного объявления
// @Tags		articles
// @Accept		json
// @Produce	json
// @Param		id	path		int	true	"ID объявления"
// @Success	200	{object}	models.Article
// @Failure	400	{string}	string	"Barticle Request"
// @Failure	500	{string}	string	"Internal Server Error"
// @Router		/article/{id} [get]
func (h *articleController) GetOneArticle(c *gin.Context) {
	go metrics.GetArticleRequest.Inc()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorLog("unable to parse id", err)
		utils.ErrorResponse(c, http.StatusBadRequest, types.ErrInvalidArticleId)
		return
	}

	article, err := h.articleService.GetOneArticle(id)
	if err != nil {
		utils.ErrorLog("unable to get article", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, article)
}
