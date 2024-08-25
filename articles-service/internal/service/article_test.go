package service

import (
	"articles-service/internal/lib/types"
	"articles-service/internal/models"
	"articles-service/internal/repository"
	"articles-service/internal/repository/mocks"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestCreateArticle(t *testing.T) {
	testArticle := &models.Article{
		Title:       "title",
		UserId:      3,
		Photos:      []string{"photo1", "photo2"},
		Description: "description",
	}
	createdArticleId := 2

	testCases := []struct {
		name    string
		mockExp func(articleRepo *mocks.ArticleRepoInterface)
		expId   int
		expErr  error
	}{
		{
			name: "Успешный кейс, все как должно работать",
			mockExp: func(articleRepo *mocks.ArticleRepoInterface) {
				articleRepo.On("Create", testArticle).
					Return(createdArticleId, nil).Times(1)
			},
			expId:  createdArticleId,
			expErr: nil,
		},
		{
			name: "База вернула любую ошибку",
			mockExp: func(articleRepo *mocks.ArticleRepoInterface) {
				articleRepo.On("Create", testArticle).
					Return(-1, errors.New("some error")).Times(1)
			},
			expId:  -1,
			expErr: errors.New("some error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockArticleRepoInterface := &mocks.ArticleRepoInterface{}
			repository := &repository.Repository{Article: mockArticleRepoInterface}
			service := NewService(repository)

			testCase.mockExp(mockArticleRepoInterface)
			id, err := service.Article.CreateArticle(testArticle)

			assert.Equal(t, testCase.expErr, err)
			assert.Equal(t, testCase.expId, id)
			mockArticleRepoInterface.AssertExpectations(t)
		})
	}
}

func TestGetAllArticles(t *testing.T) {
	var testArticles []*models.Article = []*models.Article{
		{
			Id:          1,
			Title:       "title",
			Photos:      []string{"photo1", "photo2"},
			Description: "description",
		},
		{
			Id:          2,
			Title:       "title",
			Photos:      []string{"photo1", "photo2"},
			Description: "description",
		},
	}
	testDateSort := "desc"
	testUserId := 6
	testPageNumber := 2

	testCases := []struct {
		name        string
		mockExp     func(articleRepo *mocks.ArticleRepoInterface)
		expArticles []*models.Article
		expErr      error
	}{
		{
			name: "Успешный кейс, все по плану",
			mockExp: func(articleRepo *mocks.ArticleRepoInterface) {
				articleRepo.On("GetAll", testDateSort, testPageNumber, testUserId).
					Return(testArticles, nil).Times(1)
			},
			expArticles: testArticles,
			expErr:      nil,
		},
		{
			name: "Получаем любую ошибку из базы",
			mockExp: func(articleRepo *mocks.ArticleRepoInterface) {
				articleRepo.On("GetAll", testDateSort, testPageNumber, testUserId).
					Return(nil, errors.New("some error")).Times(1)
			},
			expArticles: nil,
			expErr:      errors.New("some error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockArticleRepoInterface := &mocks.ArticleRepoInterface{}
			repository := &repository.Repository{Article: mockArticleRepoInterface}
			service := NewService(repository)

			testCase.mockExp(mockArticleRepoInterface)
			articles, err := service.Article.GetAllArticles(
				testDateSort, testPageNumber, testUserId,
			)

			assert.Equal(t, testCase.expArticles, articles)
			assert.Equal(t, testCase.expErr, err)
			mockArticleRepoInterface.AssertExpectations(t)
		})
	}
}

func TestGetOneArticle(t *testing.T) {
	testArticle := &models.Article{
		Id:          1,
		Title:       "title",
		Photos:      []string{"photo1", "photo2"},
		Description: "description",
	}
	testId := 1

	testCases := []struct {
		name       string
		mockExp    func(articleRepo *mocks.ArticleRepoInterface)
		expArticle *models.Article
		expErr     error
	}{
		{
			name: "Успешный кейс, все по плану",
			mockExp: func(articleRepo *mocks.ArticleRepoInterface) {
				articleRepo.On("GetOne", testId).
					Return(testArticle, nil).Times(1)
			},
			expArticle: testArticle,
			expErr:     nil,
		},
		{
			name: "Вернулась ошибка, что статья не найдена",
			mockExp: func(articleRepo *mocks.ArticleRepoInterface) {
				articleRepo.On("GetOne", testId).
					Return(nil, pgx.ErrNoRows).Times(1)
			},
			expArticle: nil,
			expErr:     types.ErrArticleNotFound,
		},
		{
			name: "Вернулась любая другая ошибка",
			mockExp: func(articleRepo *mocks.ArticleRepoInterface) {
				articleRepo.On("GetOne", testId).
					Return(nil, errors.New("some error")).Times(1)
			},
			expArticle: nil,
			expErr:     errors.New("some error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockArticleRepoInterface := &mocks.ArticleRepoInterface{}
			repository := &repository.Repository{Article: mockArticleRepoInterface}
			service := NewService(repository)

			testCase.mockExp(mockArticleRepoInterface)
			article, err := service.Article.GetOneArticle(testId)

			assert.Equal(t, testCase.expArticle, article)
			assert.Equal(t, testCase.expErr, err)
			mockArticleRepoInterface.AssertExpectations(t)
		})
	}
}
