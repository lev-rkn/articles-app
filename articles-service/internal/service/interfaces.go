package service

import "articles-service/internal/models"

//go:generate mockery --name ArticleServiceInterface --output ./mocks
type ArticleServiceInterface interface {
	Create(ad *models.Article) (int, error)
	GetOne(id int) (*models.Article, error)
	GetAll(dateSort string, page int, userId int) ([]*models.Article, error)
}

//go:generate mockery --name CommentServiceInterface --output ./mocks
type CommentServiceInterface interface {
	Create(comment *models.Comment) (int, error)
	GetCommentsOnArticle(articleId int) ([]*models.Comment, error)
}
