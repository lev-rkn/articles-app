package repository

import (
	"ads-service/internal/models"
	"context"
	"strings"
)

type AdRepo struct {
	ctx  context.Context
	conn PgConn
}
var _ AdRepoInterface = (*AdRepo)(nil)

func NewAdRepo(ctx context.Context,conn PgConn) *AdRepo {
	return &AdRepo{
		ctx:  context.Background(),
		conn: conn,
	}
}

func (s *AdRepo) Create(ad *models.Ad) (int, error) {
	var id int
	err := s.conn.QueryRow(s.ctx,
		"INSERT INTO advertisements (title, description, price, photos) VALUES ($1, $2, $3, $4) RETURNING id;",
		ad.Title, ad.Description, ad.Price, ad.Photos).Scan(&id)

	return id, err
}

func (s *AdRepo) GetOne(id int) (*models.Ad, error) {
	ad := &models.Ad{}

	err := s.conn.QueryRow(s.ctx,
		"SELECT title, price, photos, description FROM advertisements WHERE id = $1;", id).
		Scan(&ad.Title, &ad.Price, &ad.Photos, &ad.Description)

	return ad, err
}

func (s *AdRepo) GetAll(priceSort string, dateSort string, page int) ([]*models.Ad, error) {
	// сборка строки сортировки по цене и по дате для запроса
	var orderQuery string
	sorts := make([]string, 0, 2)
	if priceSort != "" {
		sorts = append(sorts, priceSort)
	}
	if dateSort != "" {
		sorts = append(sorts, dateSort)
	}
	if len(sorts) > 0 {
		orderQuery += " ORDER BY " + strings.Join(sorts, ", ")
	}

	const adsPerPage = 10
	var skipAds = adsPerPage * (page)

	rows, err := s.conn.Query(s.ctx,
		"SELECT id, title, price, photos FROM advertisements WHERE id > $1"+orderQuery+" LIMIT $2;",
		skipAds, adsPerPage)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Перекладывание всех объектов объявлений в массив
	adsArr := make([]*models.Ad, 0, adsPerPage)
	for rows.Next() {
		ad := &models.Ad{}
		err = rows.Scan(&ad.Id, &ad.Title, &ad.Price, &ad.Photos)
		if err != nil {
			return nil, err
		}
		adsArr = append(adsArr, ad)
	}

	return adsArr, nil
}
