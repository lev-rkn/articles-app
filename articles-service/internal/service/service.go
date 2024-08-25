package service

import (
	"articles-service/internal/repository"
)

type Service struct {
	Article ArticleServiceInterface
	Comment CommentServiceInterface
}

func NewService(
	repository *repository.Repository,
) *Service {
	return &Service{
		Article: &articleService{repository: repository},
		Comment: &commentService{repository: repository},
	}
}
