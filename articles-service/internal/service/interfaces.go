package service

import "articles-service/internal/models"

//go:generate mockery --name ArticleServiceInterface --output ./mocks
type ArticleServiceInterface interface {
	Create(ad *models.Article) (int, error)
	GetOne(id int) (*models.Article, error)
	GetAll(priceSort string, dateSort string, page int, userId int) ([]*models.Article, error)
}

// TODO comment interface