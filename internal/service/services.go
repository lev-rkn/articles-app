package services

import (
	"ads-service/internal/models"
	"ads-service/internal/repository"
	"log/slog"
)

type AdServiceInterface interface {
	Create(ad *models.Ad) (int, error)
	GetOne(id int) (*models.Ad, error)
	GetAll(priceSort string, dateSort string, page int) ([]*models.Ad, error)
}

type Service struct {
	Ad AdServiceInterface
}

func NewService(
	repository *repository.Repository,
	logger *slog.Logger,
) *Service {
	return &Service{
		Ad: &adService{repository: repository, logger: logger},
	}
}
