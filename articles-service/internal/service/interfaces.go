package service

import "articles-service/internal/models"

//go:generate mockery --name ArticleServiceInterface --output ./mocks
type ArticleServiceInterface interface {
	CreateArticle(ad *models.Article) (int, error)
	GetOneArticle(id int) (*models.Article, error)
	GetAllArticles(dateSort string, page int, userId int) ([]*models.Article, error)
}

//go:generate mockery --name CommentServiceInterface --output ./mocks
type CommentServiceInterface interface {
	CreateComment(comment *models.Comment) (int, error)
	GetCommentsOnArticle(articleId int) ([]*models.Comment, error)
}
