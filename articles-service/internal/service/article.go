package service

import (
	"articles-service/internal/models"
	"articles-service/internal/repository"
	"articles-service/metrics"
	"fmt"
)

type articleService struct {
	repository *repository.Repository
}

var _ ArticleServiceInterface = (*articleService)(nil)

func (s *articleService) CreateArticle(article *models.Article) (int, error) {
	var err error
	defer func() {
		if err == nil {
			go metrics.CreateArticleOK.Inc()
		} else {
			go metrics.CreateArticleError.Inc()
		}
	}()

	id, err := s.repository.Article.Create(article)
	if err != nil {
		return -1, fmt.Errorf("repository.Article.Create: %w", err)
	}

	return id, nil
}

func (s *articleService) GetAllArticles(dateSort string, page int, userId int,
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
	if err != nil {
		return nil, fmt.Errorf("repository.Article.GetAll: %w", err)
	}

	return articles, nil
}

func (s *articleService) GetOneArticle(id int) (*models.Article, error) {
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
		return nil, fmt.Errorf("repository.Article.GetOne: %w", err)
	}

	return article, nil
}
