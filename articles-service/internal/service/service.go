package service

import (
	"articles-service/internal/repository"
)



type Service struct {
	Article ArticleServiceInterface
}

func NewService(
	repository *repository.Repository,
) *Service {
	return &Service{
		Article: &articleService{repository: repository},
	}
}
