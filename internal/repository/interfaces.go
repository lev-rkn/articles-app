package repository

import (
	"ads-service/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name AdRepoInterface --output ./mocks
type AdRepoInterface interface {
	Create(ad *models.Ad) (int, error)
	GetOne(id int) (*models.Ad, error)
	GetAll(priceSort string, dateSort string, page int) ([]*models.Ad, error)
}

//go:generate mockery --name PgConn  --output ./mocks
type PgConn interface {
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
}