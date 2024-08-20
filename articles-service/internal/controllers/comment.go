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

type commentController struct {
	ctx            context.Context
	commentService service.CommentServiceInterface
}

func InitCommentController(
	ctx context.Context,
	commentService service.CommentServiceInterface,
	router *gin.RouterGroup,
) *commentController {
	commentController := &commentController{
		ctx:            ctx,
		commentService: commentService,
	}

	router.POST("/create/", commentController.CreateComment)
	router.GET("/:articleId", commentController.GetCommentsOnArticle)

	return commentController
}

//	@Summary	Создание комментария
//	@Tags		comments
//	@Accept		json
//	@Produce	json
//	@Param		comment	body		models.Comment	true	"Комментарий"
//	@Success	201		{int}		id
//	@Failure	400		{string}	string	"Bad Request"
//	@Failure	500		{string}	string	"Internal Server Error"
//	@Router		/comments/create/ [post]
func (h *commentController) CreateComment(c *gin.Context) {
	go metrics.CreateCommentRequest.Inc()
	// проверяем наличие ошибки, возможно переданной нам через middleware
	if err, ok := c.Value(types.KeyError).(error); ok {
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, err)
			return
		}
	}

	comment := &models.Comment{}
	err := json.NewDecoder(c.Request.Body).Decode(&comment)
	if err != nil {
		utils.ErrorLog("unable to decode comment", err)
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	parsedClaims := utils.ParseUserClaims(c.Value(types.KeyUser))
	comment.UserId = parsedClaims.UserId

	validate := validator.New()
	err = validate.Struct(comment)
	if err != nil {
		utils.ErrorLog("validate comment", err)
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.commentService.CreateComment(comment)
	if err != nil {
		utils.ErrorLog("unable to create comment", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

//	@Summary	Получение всех комментариев статьи
//	@Tags		comments
//	@Accept		json
//	@Produce	json
//	@Success	200	{array}		models.Comment
//	@Failure	400	{string}	string			"Bad Request"
//	@Failure	500	{string}	string			"Internal Server Error"
//	@Router		/comments/{articleId} [get]
func (h *commentController) GetCommentsOnArticle(c *gin.Context) {
	go metrics.GetCommentsRequest.Inc()

	articleId, err := strconv.Atoi(c.Param("articleId"))
	if err != nil {
		utils.ErrorLog("parse articleId", err)
		utils.ErrorResponse(c, http.StatusBadRequest, types.ErrInvalidArticleId)
		return
	}

	commentsArr, err := h.commentService.GetCommentsOnArticle(articleId)
	if err != nil {
		utils.ErrorLog("unable to get comments", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, commentsArr)
}
