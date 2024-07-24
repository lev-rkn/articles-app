package services

import (
	"ads-service/internal/models"
	"ads-service/internal/store"
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
	store *store.Store,
	logger *slog.Logger,
) *Service {
	return &Service{
		Ad: &adService{store: store, logger: logger},
	}
}
