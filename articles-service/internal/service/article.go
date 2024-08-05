package service

import (
	"articles-service/internal/lib/types"
	"articles-service/internal/lib/utils"
	"articles-service/internal/models"
	"articles-service/internal/repository"
	"errors"

	"github.com/jackc/pgx/v5"
)

type articleService struct {
	repository *repository.Repository
}

var _ ArticleServiceInterface = (*articleService)(nil)

func (s *articleService) GetAll(priceSort string, dateSort string, page int, userId int,
) ([]*models.Article, error) {
	articles, err := s.repository.Article.GetAll(priceSort, dateSort, page, userId)
	// TODO: никакой обработки ошибок из базы
	if err != nil {
		utils.ErrorLog("service.article.GetAll", err)
		return nil, err
	}

	return articles, nil
}

func (s *articleService) Create(article *models.Article) (int, error) {
	id, err := s.repository.Article.Create(article)
	// TODO: никакой обработки ошибок из базы
	if err != nil {
		utils.ErrorLog("service.article.Create", err)
		return -1, err
	}

	return id, nil
}

func (s *articleService) GetOne(id int) (*models.Article, error) {
	article, err := s.repository.Article.GetOne(id)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, types.ErrArticleNotFound
	}

	return article, nil
}
