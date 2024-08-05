package controllers

import (
	"articles-service/internal/controllers/middlewares"
	"articles-service/internal/lib/types"
	"articles-service/internal/models"
	"articles-service/internal/service"
	"articles-service/metrics"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
)

type commentController struct {
	ctx            context.Context
	commentService service.CommentServiceInterface
	router         *gin.RouterGroup
}

func InitCommentController(
	ctx context.Context,
	commentService service.CommentServiceInterface,
	router *gin.RouterGroup,
) *commentController {
	commentController := &commentController{
		ctx:            ctx,
		commentService: commentService,
		router:         router,
	}
	router.GET("/:articleId/", commentController.GetCommentsOnArticle)
	router.Use(middlewares.AuthMiddleware())
	router.POST("/create/", commentController.CreateComment)

	return commentController
}

// @Summary Создание комментария
// @Tags comments
// @Accept json
// @Produce json
// @Param comment body models.Comment true "Комментарий"
// @Success 201 {int} id
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router			/comments/create/ [post]
func (h *commentController) CreateComment(c *gin.Context) {
	go metrics.CreateCommentRequest.Inc()
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
			go metrics.CreateCommentOK.Inc()
		} else {
			go metrics.CreateCommentError.Inc()
		}
	}()

	comment := &models.Comment{}
	err = json.NewDecoder(c.Request.Body).Decode(&comment)
	if err != nil {
		slog.Error("unable to decode comment", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// достаем токен из контекста и берем оттуда id пользователя, создавшего этот коммент
	if v, ok := c.Value(types.KeyUser).(*jwt.Token); ok {
		if claims, ok := v.Claims.(jwt.MapClaims); ok {
			if idF, ok := claims["uid"].(float64); ok {
				id := int(idF)
				comment.UserId = id
			}
		}
	}

	validate := validator.New()
	err = validate.Struct(comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.commentService.Create(comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// @Summary Получение всех комментариев статьи
// @Tags comments
// @Accept json
// @Produce json
// @Success 200 {array} models.Comment true "Комментрии"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router			/comments/{articleId}/ [get]
func (h *commentController) GetCommentsOnArticle(c *gin.Context) {
	var err error
	// собираем метрики
	go metrics.GetCommentsRequest.Inc()
	defer func() {
		if err == nil {
			go metrics.GetCommentsOK.Inc()
		} else {
			go metrics.GetCommentsError.Inc()
		}
	}()

	articleId, err := strconv.Atoi(c.Param("articleId"))
	if err != nil {
		slog.Error("parse articleId", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	commentsArr, err := h.commentService.GetCommentsOnArticle(articleId)
	if err != nil {
		slog.Error("unable to get comments", "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, commentsArr)
}
