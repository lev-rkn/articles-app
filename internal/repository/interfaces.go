package repository

import "ads-service/internal/models"

//go:generate mockery --name AdRepo --output ./mocks
type AdRepo interface {
	Create(ad *models.Ad) (int, error)
	GetOne(id int) (*models.Ad, error)
	GetAll(priceSort string, dateSort string, page int) ([]*models.Ad, error)
}
