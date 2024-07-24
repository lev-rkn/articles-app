package services

import (
	"ads-service/internal/models"
	"ads-service/internal/store"
	"log/slog"
)
type adService struct {
	store *store.Store
	logger *slog.Logger
}
var _ AdServiceInterface = (*adService)(nil)

func (s *adService) GetAll(priceSort string, dateSort string, page int) ([]*models.Ad, error) {
	ads, err := s.store.Ad.GetAll(priceSort, dateSort, page)
	if err != nil {
		s.logger.Error("unable to get ads", "err", err.Error())
		return nil, err
	}

	return ads, nil
}

func (s *adService) Create(ad *models.Ad) (int, error) {
	int, err := s.store.Ad.Create(ad)
	if err != nil {
		s.logger.Error("unable to create ad", "err", err.Error())
		return 0, err
	}

	return int, nil
}

func (s *adService) GetOne(id int) (*models.Ad, error) {
	ad, err := s.store.Ad.GetOne(id)
	if err != nil {
		s.logger.Error("unable to get ad", "err", err.Error())
		return nil, err
	}

	return ad, nil
}
