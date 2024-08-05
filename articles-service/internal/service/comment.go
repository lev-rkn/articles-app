package service

import (
	"articles-service/internal/lib/types"
	"articles-service/internal/models"
	"articles-service/internal/repository"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type commentService struct {
	repository *repository.Repository
}

var _ CommentServiceInterface = (*commentService)(nil)

func (s *commentService) GetCommentsOnArticle(articleId int) ([]*models.Comment, error) {
	comments, err := s.repository.Comment.GetCommentsOnArticle(articleId)
	if err != nil {
		slog.Error("service.comment.GetAll", "err", err.Error())
		return nil, err
	}

	return comments, nil
}

func (s *commentService) Create(comment *models.Comment) (int, error) {
	// проверяем, что статья действительно существует
	_, err := s.repository.Article.GetOne(comment.ArticleId)
	if errors.Is(err, pgx.ErrNoRows) {
		return -1, types.ErrArticleNotFound
	}
	if err != nil {
		slog.Error("service.comment.create", "err", err.Error())
		return -1, err
	}

	id, err := s.repository.Comment.Create(comment)
	if err != nil {
		slog.Error("service.comment.Create", "err", err.Error())
		return -1, err
	}

	return id, nil
}
