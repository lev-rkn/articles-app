package service

import (
	"articles-service/internal/models"
	"articles-service/internal/repository"
	"articles-service/metrics"
	"fmt"
)

type commentService struct {
	repository *repository.Repository
}

var _ CommentServiceInterface = (*commentService)(nil)

func (s *commentService) CreateComment(comment *models.Comment) (int, error) {
	var err error
	defer func() {
		if err == nil {
			go metrics.CreateCommentOK.Inc()
		} else {
			go metrics.CreateCommentError.Inc()
		}
	}()

	// проверяем, что статья действительно существует
	_, err = s.repository.Article.GetOne(comment.ArticleId)
	if err != nil {
		return -1, fmt.Errorf("repository.Article.GetOne %w", err)
	}

	id, err := s.repository.Comment.Create(comment)
	if err != nil {
		return -1, fmt.Errorf("repository.Comment.Create: %w", err)
	}

	return id, nil
}

func (s *commentService) GetCommentsOnArticle(articleId int) ([]*models.Comment, error) {
	var err error
	defer func() {
		if err == nil {
			go metrics.GetCommentsOK.Inc()
		} else {
			go metrics.GetCommentsError.Inc()
		}
	}()

	comments, err := s.repository.Comment.GetCommentsOnArticle(articleId)
	if err != nil {
		return nil, fmt.Errorf("repository.Comment.GetCommentsOnArticle: %w", err)
	}

	return comments, nil
}
