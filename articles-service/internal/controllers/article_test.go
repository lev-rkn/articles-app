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

func TestCreateArticle(t *testing.T) {
	testArticle := &models.Article{
		Title:       "title",
		Description: "fnsfbn",
		Photos:      []string{"photo1", "photo2"},
	}
	marshalledArticle, _ := json.Marshal(testArticle)
	articleInJSON := string(marshalledArticle)
	testArticleForValidationTest, _ := json.Marshal(&models.Article{
		Title:       "title",
		Description: "",
		Photos:      []string{"photo1", "photo2"},
	})

	testCases := []struct {
		name        string
		mockExpect  func(articleService *mocks.ArticleServiceInterface)
		reqData     string         // то, что мы кидаем в нашу тестируему функцию как тело запроса
		inContext   map[string]any // внедряем в контекст роутера
		expJSON     string
		expError    bool // если ожидается ошибка - мы точно не знаем какая
		expHTTPCode int
	}{
		{
			name: "Успешный кейс, все работает как и задумывалось",
			mockExpect: func(articleService *mocks.ArticleServiceInterface) {
				articleService.On("CreateArticle", testArticle).Return(3, nil)
			},
			reqData:     articleInJSON,
			expJSON:     `{"id":3}`,
			expHTTPCode: http.StatusCreated,
		},
		{
			name:        "Подсовываем невалидный JSON",
			mockExpect:  func(articleService *mocks.ArticleServiceInterface) {},
			reqData:     `{""""""""""""""}`,
			expError:    true,
			expHTTPCode: http.StatusBadRequest,
		},
		{
			name:       "Подсовываем ошибку в контексте",
			mockExpect: func(articleService *mocks.ArticleServiceInterface) {},
			reqData:    articleInJSON,
			inContext: map[string]any{
				types.KeyError: errors.New("some error"),
			},
			expError:    true,
			expHTTPCode: http.StatusBadRequest,
		},
		{
			name:        "Проверяем, что валидация вообще работает",
			mockExpect:  func(articleService *mocks.ArticleServiceInterface) {},
			reqData:     string(testArticleForValidationTest),
			expError:    true,
			expHTTPCode: http.StatusBadRequest,
		},
		{
			name: "Любая ошибка от сервиса статей",
			mockExpect: func(articleService *mocks.ArticleServiceInterface) {
				articleService.On("CreateArticle", testArticle).
					Return(-1, errors.New("some error"))
			},
			reqData:     articleInJSON,
			expError:    true,
			expHTTPCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Иницализируем все зависимости
			ctx := context.Background()
			mockArticleService := &mocks.ArticleServiceInterface{}
			router := gin.Default()
			articleRouter := router.Group("/article")
			// внедряем нужные нам для теста значения в контекст через middleware
			// перед инициализацией тестируемого контроллера
			articleRouter.Use(func() gin.HandlerFunc {
				return func(c *gin.Context) {
					for k, v := range testCase.inContext {
						c.Set(k, v)
					}
				}
			}())
			InitArticleController(ctx, mockArticleService, articleRouter)

			// подготовка к запросу
			testCase.mockExpect(mockArticleService)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(
				"POST", "/article/create/",
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
			mockArticleService.AssertExpectations(t)
		})
	}
}

func TestGetAllArticles(t *testing.T) {
	testArticles := []*models.Article{
		{
			Title:       "title",
			Description: "fnsfbn",
			Photos:      []string{"photo1", "photo2"},
		},
		{
			Title:       "title",
			Description: "fnsfbn",
			Photos:      []string{"photo1", "photo2"},
		},
	}
	marshalledArticles, _ := json.Marshal(testArticles)
	articlesInJSON := string(marshalledArticles)

	testCases := []struct {
		name        string
		mockExp     func(articleService *mocks.ArticleServiceInterface)
		queryParams map[string]string
		expJSON     string
		expHTTPCode int
	}{
		{
			name: "Успешный кейс, рядовая ситуация",
			mockExp: func(articleService *mocks.ArticleServiceInterface) {
				articleService.On("GetAllArticles", "asc", 4, 7).Return(testArticles, nil)
			},
			queryParams: map[string]string{"page": "4", "user_id": "7", "date": "asc"},
			expJSON:     articlesInJSON,
			expHTTPCode: http.StatusOK,
		},
		{
			name: "Любая ошибка из сервиса статей",
			mockExp: func(articleService *mocks.ArticleServiceInterface) {
				articleService.On("GetAllArticles", "asc", 4, 7).
					Return(nil, errors.New("some error"))
			},
			queryParams: map[string]string{"page": "4", "user_id": "7", "date": "asc"},
			expJSON:     `{"error":"some error"}`,
			expHTTPCode: http.StatusInternalServerError,
		},
		{
			name:        "Неверный номер страницы",
			mockExp:     func(articleService *mocks.ArticleServiceInterface) {},
			queryParams: map[string]string{"user_id": "7", "date": "asc"},
			expJSON:     fmt.Sprintf(`{"error":"%s"}`, types.ErrInvalidPageNumber.Error()),
			expHTTPCode: http.StatusBadRequest,
		},
		{
			name:        "Неверный идентификатор пользователя",
			mockExp:     func(articleService *mocks.ArticleServiceInterface) {},
			queryParams: map[string]string{"page": "4", "user_id": "7п", "date": "asc"},
			expJSON:     fmt.Sprintf(`{"error":"%s"}`, types.ErrInvalidUserId.Error()),
			expHTTPCode: http.StatusBadRequest,
		},
		{
			name:        "Неверный параметр сортировки по дате",
			mockExp:     func(articleService *mocks.ArticleServiceInterface) {},
			queryParams: map[string]string{"page": "4", "date": "ascending)"},
			expJSON:     fmt.Sprintf(`{"error":"%s"}`, types.ErrInvalidDateSort.Error()),
			expHTTPCode: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()
			mockArticleService := &mocks.ArticleServiceInterface{}
			router := gin.Default()
			articleRouter := router.Group("/article")
			InitArticleController(ctx, mockArticleService, articleRouter)

			testCase.mockExp(mockArticleService)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/article/all/", strings.NewReader(""))

			query := req.URL.Query()
			for k, v := range testCase.queryParams {
				query.Add(k, v)
			}
			req.URL.RawQuery = query.Encode()
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expJSON, w.Body.String())
			assert.Equal(t, testCase.expHTTPCode, w.Code)
			mockArticleService.AssertExpectations(t)
		})
	}
}

func TestGetOneArticle(t *testing.T) {
	testArticle := &models.Article{
		Title:       "title",
		Description: "fnsfbn",
		Photos:      []string{"photo1", "photo2"},
	}
	marshalledArticle, _ := json.Marshal(testArticle)
	articleInJSON := string(marshalledArticle)

	testCases := []struct {
		name        string
		mockExp     func(articleService *mocks.ArticleServiceInterface)
		articleId   string
		expJSON     string
		expHTTPCode int
	}{
		{
			name: "Успешный кейс, рядовая ситуация",
			mockExp: func(articleService *mocks.ArticleServiceInterface) {
				articleService.On("GetOneArticle", 4).Return(testArticle, nil)
			},
			articleId:   "4",
			expJSON:     articleInJSON,
			expHTTPCode: http.StatusOK,
		},
		{
			name:        "Невалидный идентификатор страницы",
			mockExp:     func(articleService *mocks.ArticleServiceInterface) {},
			articleId:   "4a",
			expJSON:     fmt.Sprintf(`{"error":"%s"}`, types.ErrInvalidArticleId.Error()),
			expHTTPCode: http.StatusBadRequest,
		},
		{
			name: "Любая ошибка из сервиса объявлений",
			mockExp: func(articleService *mocks.ArticleServiceInterface) {
				articleService.On("GetOneArticle", 4).
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
			mockArticleService := &mocks.ArticleServiceInterface{}
			router := gin.Default()
			articleRouter := router.Group("/article")
			InitArticleController(ctx, mockArticleService, articleRouter)

			testCase.mockExp(mockArticleService)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(
				"GET", "/article/"+testCase.articleId,
				strings.NewReader(""))
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expJSON, w.Body.String())
			assert.Equal(t, testCase.expHTTPCode, w.Code)
			mockArticleService.AssertExpectations(t)
		})
	}
}
