package services

import (
	"ads-service/internal/models"
	"ads-service/internal/repository"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type adService struct {
	repository *repository.Repository
	logger     *slog.Logger
}

var _ AdServiceInterface = (*adService)(nil)

// Ошибки
var (
	errAdNotFound = errors.New("ad not found")
)

func (s *adService) GetAll(priceSort string, dateSort string, page int) ([]*models.Ad, error) {
	ads, err := s.repository.Ad.GetAll(priceSort, dateSort, page)
	if err != nil {
		slog.Error("unable to get ads", "err", err.Error())
		return nil, err
	}

	return ads, nil
}

func (s *adService) Create(ad *models.Ad) (int, error) {
	int, err := s.repository.Ad.Create(ad)
	if err != nil {
		slog.Error("unable to create ad", "err", err.Error())
		return -1, err
	}

	return int, nil
}

func (s *adService) GetOne(id int) (*models.Ad, error) {
	ad, err := s.repository.Ad.GetOne(id)

	if err == pgx.ErrNoRows {
		return nil, errAdNotFound
	}

	if err != nil {
		slog.Error("unable to get ad", "err", err.Error())
		return nil, err
	}

	return ad, nil
}
