package service

import (
	"articles-service/internal/lib/types"
	"articles-service/internal/lib/utils"
	"articles-service/internal/models"
	"articles-service/internal/repository"
	"errors"

	"github.com/jackc/pgx/v5"
)

type commentService struct {
	repository *repository.Repository
}

var _ CommentServiceInterface = (*commentService)(nil)

func (s *commentService) GetCommentsOnArticle(articleId int) ([]*models.Comment, error) {
	comments, err := s.repository.Comment.GetCommentsOnArticle(articleId)
	if err != nil {
		utils.ErrorLog("service.comment.GetAll", err)
		return nil, err
	}

	return comments, nil
}

func (s *commentService) Create(comment *models.Comment) (int, error) {
	// проверяем, что статья действительно существует
	_, err := s.repository.Article.GetOne(comment.ArticleId)
	if err != nil {
		utils.ErrorLog("service.comment.create", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, types.ErrArticleNotFound
		}
		return -1, err
	}

	id, err := s.repository.Comment.Create(comment)
	// TODO: никакой обработки ошибок из базы
	if err != nil {
		utils.ErrorLog("service.comment.Create", err)
		return -1, err
	}

	return id, nil
}
