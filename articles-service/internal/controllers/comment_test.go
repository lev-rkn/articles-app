package controllers

import (
	"articles-service/internal/lib/types"
	"articles-service/internal/models"
	"articles-service/internal/service/mocks"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	testComment := &models.Comment{
		ArticleId: 4,
		UserId:    99,
		Text:      "poeirhngpeoirghj",
	}
	marshalledComment, _ := json.Marshal(testComment)
	commentInJSON := string(marshalledComment)
	testCommentForValidationTest, _ := json.Marshal(&models.Comment{
		UserId: 99,
		Text:   "poeirhngpeoirghj",
	})

	testCases := []struct {
		name        string
		mockExpect  func(commentService *mocks.CommentServiceInterface)
		reqData     string // то, что мы кидаем в нашу тестируему функцию как тело запроса
		inContext   map[string]any // внедряем в контекст роутера
		expJSON     string
		expError    bool // если ожидается ошибка - мы точно не знаем какая
		expHTTPCode int
	}{
		{
			name: "Успешный кейс, все работает как и задумывалось",
			mockExpect: func(commentService *mocks.CommentServiceInterface) {
				commentService.On("CreateComment", testComment).Return(3, nil)
			},
			reqData:     commentInJSON,
			expJSON:     `{"id":3}`,
			expHTTPCode: http.StatusCreated,
		},
		{
			name:        "Подсовываем невалидный JSON",
			mockExpect:  func(commentService *mocks.CommentServiceInterface) {},
			reqData:     `{""""""""""""""}`,
			expError:    true,
			expHTTPCode: http.StatusBadRequest,
		},
				{
			name:       "Подсовываем ошибку в контексте",
			mockExpect: func(commentService *mocks.CommentServiceInterface) {},
			reqData:    commentInJSON,
			inContext: map[string]any{
				types.KeyError: errors.New("some error"),
			},
			expError:    true,
			expHTTPCode: http.StatusBadRequest,
		},
		{
			name:        "Проверяем, что валидация вообще работает",
			mockExpect:  func(commentService *mocks.CommentServiceInterface) {},
			reqData:     string(testCommentForValidationTest),
			expError:    true,
			expHTTPCode: http.StatusBadRequest,
		},
		{
			name: "Любая ошибка от сервиса статей",
			mockExpect: func(commentService *mocks.CommentServiceInterface) {
				commentService.On("CreateComment", testComment).
					Return(-1, errors.New("some error"))
			},
			reqData:     commentInJSON,
			expError:    true,
			expHTTPCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Иницализируем все зависимости
			ctx := context.Background()
			mockCommentService := &mocks.CommentServiceInterface{}
			router := gin.Default()
			commentRouter := router.Group("/comment")
			commentRouter.Use(func() gin.HandlerFunc {
				return func(c *gin.Context) {
					for k, v := range testCase.inContext {
						c.Set(k, v)
					}
				}
			}())
			InitCommentController(ctx, mockCommentService, commentRouter)

			// подготовка к запросу
			testCase.mockExpect(mockCommentService)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(
				"POST", "/comment/create/",
				strings.NewReader(testCase.reqData),
			)
			router.ServeHTTP(w, req)

			if testCase.expError {
				// какая точно ошибка - нам наверняка неизвестно
				assert.Contains(t, w.Body.String(), `"error":`)
			} else {
				assert.Equal(t, testCase.expJSON, w.Body.String())
			}
			assert.Equal(t, testCase.expHTTPCode, w.Code)
			mockCommentService.AssertExpectations(t)
		})
	}
}

func TestGetCommentsOnArticle(t *testing.T) {
	testComments := []*models.Comment{
		{
			ArticleId: 4,
			UserId:    99,
			Text:      "poeirhngpeoirghj",
		},
		{
			ArticleId: 4,
			UserId:    99,
			Text:      "poeirhngpeoirghj",
		},
	}
	marshalledComments, _ := json.Marshal(testComments)
	commentsInJSON := string(marshalledComments)

	testCases := []struct {
		name        string
		mockExp     func(commentService *mocks.CommentServiceInterface)
		articleId   string
		expJSON     string
		expHTTPCode int
	}{
		{
			name: "Успешный кейс, получаем все комментарии беспрепятственно",
			mockExp: func(commentService *mocks.CommentServiceInterface) {
				commentService.On("GetCommentsOnArticle", 4).Return(testComments, nil)
			},
			articleId:   "4",
			expJSON:     commentsInJSON,
			expHTTPCode: http.StatusOK,
		},
		{
			name:        "Невалидный идентификатор статьи",
			mockExp:     func(commentService *mocks.CommentServiceInterface) {},
			articleId:   "hfa",
			expJSON:     fmt.Sprintf(`{"error":"%s"}`, types.ErrInvalidArticleId.Error()),
			expHTTPCode: http.StatusBadRequest,
		},
		{
			name: "Любая ошибка из сервиса комментариев",
			mockExp: func(commentService *mocks.CommentServiceInterface) {
				commentService.On("GetCommentsOnArticle", 4).
					Return(nil, errors.New("some error"))
			},
			articleId:   "4",
			expJSON:     `{"error":"some error"}`,
			expHTTPCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()
			mockCommentService := &mocks.CommentServiceInterface{}
			router := gin.Default()
			commentRouter := router.Group("/comments")
			InitCommentController(ctx, mockCommentService, commentRouter)

			testCase.mockExp(mockCommentService)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(
				"GET", "/comments/"+testCase.articleId, strings.NewReader(""),
			)
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expJSON, w.Body.String())
			assert.Equal(t, testCase.expHTTPCode, w.Code)
			mockCommentService.AssertExpectations(t)
		})
	}
}
