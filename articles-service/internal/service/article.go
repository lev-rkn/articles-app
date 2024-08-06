package service

import (
	"articles-service/internal/lib/types"
	"articles-service/internal/lib/utils"
	"articles-service/internal/models"
	"articles-service/internal/repository"
	"articles-service/metrics"
	"errors"

	"github.com/jackc/pgx/v5"
)

type articleService struct {
	repository *repository.Repository
}

var _ ArticleServiceInterface = (*articleService)(nil)

func (s *articleService) Create(article *models.Article) (int, error) {
	var err error
	defer func() {
		if err == nil {
			go metrics.CreateArticleOK.Inc()
		} else {
			go metrics.CreateArticleError.Inc()
		}
	}()

	id, err := s.repository.Article.Create(article)
	// TODO: никакой обработки ошибок из базы
	if err != nil {
		utils.ErrorLog("service.article.Create", err)
		return -1, err
	}

	return id, nil
}

func (s *articleService) GetAll(dateSort string, page int, userId int,
) ([]*models.Article, error) {
	var err error
	defer func() {
		if err == nil {
			go metrics.GetArticlesOK.Inc()
		} else {
			go metrics.GetArticlesError.Inc()
		}
	}()

	articles, err := s.repository.Article.GetAll(dateSort, page, userId)
	// TODO: никакой обработки ошибок из базы
	if err != nil {
		utils.ErrorLog("service.article.GetAll", err)
		return nil, err
	}

	return articles, nil
}

func (s *articleService) GetOne(id int) (*models.Article, error) {
	var err error
	defer func() {
		if err == nil {
			go metrics.GetArticleOK.Inc()
		} else {
			go metrics.GetArticleError.Inc()
		}
	}()

	article, err := s.repository.Article.GetOne(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, types.ErrArticleNotFound
		}
		return nil, err
	}

	return article, nil
}
