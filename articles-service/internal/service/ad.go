package service

import (
	"articles-service/internal/lib/types"
	"articles-service/internal/models"
	"articles-service/internal/repository"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type articleService struct {
	repository *repository.Repository
}

var _ ArticleServiceInterface = (*articleService)(nil)

func (s *articleService) GetAll(priceSort string, dateSort string, page int, userId int,
	) ([]*models.Article, error) {
	articles, err := s.repository.Article.GetAll(priceSort, dateSort, page, userId)
	if err != nil {
		slog.Error("service.article.GetAll", "err", err.Error())
		return nil, err
	}

	return articles, nil
}

func (s *articleService) Create(article *models.Article) (int, error) {
	id, err := s.repository.Article.Create(article)
	if err != nil {
		slog.Error("service.article.Create", "err", err.Error())
		return -1, err
	}

	return id, nil
}

func (s *articleService) GetOne(id int) (*models.Article, error) {
	article, err := s.repository.Article.GetOne(id)

	if err == pgx.ErrNoRows {
		return nil, types.ErrArticleNotFound
	}

	if err != nil {
		slog.Error("service.article.GetOne", "err", err.Error())
		return nil, err
	}

	return article, nil
}
