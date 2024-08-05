package repository

import (
	"ads-service/internal/models"
	"context"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/georgysavva/scany/v2/pgxscan"
)

type AdRepo struct {
	ctx  context.Context
	conn *pgx.Conn
}

var _ AdRepoInterface = (*AdRepo)(nil)

func NewAdRepo(ctx context.Context, conn *pgx.Conn) *AdRepo {
	return &AdRepo{
		ctx:  ctx,
		conn: conn,
	}
}

func (s *AdRepo) Create(ad *models.Ad) (int, error) {
	var id int
	err := s.conn.QueryRow(s.ctx,
		`INSERT INTO advertisements (title, description, price, photos, user_id) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
		ad.Title, ad.Description, ad.Price, ad.Photos, ad.UserId).Scan(&id)

	return id, err
}

func (s *AdRepo) GetOne(id int) (*models.Ad, error) {
	ad := &models.Ad{}
	err := pgxscan.Get(
		s.ctx, s.conn, ad, `
		SELECT *
		FROM advertisements 
		WHERE id=$1;`, id,
	)

	return ad, err
}

func (s *AdRepo) GetAll(priceSort string, dateSort string, page int, userId int) ([]*models.Ad, error) {
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

	var userIdFilterQuery = ""
	if userId > 0 {
		id := strconv.Itoa(userId)
		userIdFilterQuery = "WHERE user_id = " + id
	}

	const adsPerPage = 10
	var skipAds = adsPerPage * (page)

	rows, err := s.conn.Query(s.ctx,
		`SELECT *
		FROM advertisements
		`+userIdFilterQuery+orderQuery+` LIMIT $1 OFFSET $2;`,
		adsPerPage, skipAds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Перекладывание всех объектов объявлений в массив
	adsArr := make([]*models.Ad, 0, adsPerPage)
	err = pgxscan.ScanAll(&adsArr, rows)
	if err != nil {
		return nil, err
	}

	return adsArr, nil
}
