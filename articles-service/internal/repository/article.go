package repository

import (
	"articles-service/internal/models"
	"context"
	"strconv"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type ArticleRepo struct {
	ctx  context.Context
	conn *pgx.Conn
}

var _ ArticleRepoInterface = (*ArticleRepo)(nil)

func NewArticleRepo(ctx context.Context, conn *pgx.Conn) *ArticleRepo {
	return &ArticleRepo{
		ctx:  ctx,
		conn: conn,
	}
}

func (s *ArticleRepo) Create(article *models.Article) (int, error) {
	var id int
	err := s.conn.QueryRow(s.ctx,
		`INSERT INTO articles (title, description, photos, user_id) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id;`,
		article.Title, 
		article.Description, 
		article.Photos, 
		article.UserId).Scan(&id)

	return id, err
}

func (s *ArticleRepo) GetOne(id int) (*models.Article, error) {
	article := &models.Article{}
	err := pgxscan.Get(
		s.ctx, s.conn, article, `
		SELECT *
		FROM articles 
		WHERE id=$1;`, id,
	)

	return article, err
}

func (s *ArticleRepo) GetAll(priceSort string, dateSort string, page int, userId int,
) ([]*models.Article, error) {
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

	const articlesPerPage = 10
	var skipArticles = articlesPerPage * (page)

	rows, err := s.conn.Query(s.ctx,
		`SELECT *
		FROM articles
		`+userIdFilterQuery+orderQuery+` LIMIT $1 OFFSET $2;`,
		articlesPerPage, skipArticles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Перекладывание всех объектов объявлений в массив
	articlesArr := make([]*models.Article, 0, articlesPerPage)
	err = pgxscan.ScanAll(&articlesArr, rows)
	if err != nil {
		return nil, err
	}

	return articlesArr, nil
}
