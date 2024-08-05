package repository

import "articles-service/internal/models"

//go:generate mockery --name ArticleRepoInterface --output ./mocks
type ArticleRepoInterface interface {
	Create(ad *models.Article) (int, error)
	GetOne(id int) (*models.Article, error)
	GetAll(priceSort string, dateSort string, page int, userId int) ([]*models.Article, error)
}

// TODO comment interface
