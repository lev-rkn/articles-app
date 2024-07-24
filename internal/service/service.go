package service

import (
	"ads-service/internal/models"
)

type AdService interface {
	Create(ad *models.Ad) (int, error)
	GetOne(id int) (*models.Ad, error)
	GetAll(priceSort string, dateSort string, page int) ([]*models.Ad, error)
}

