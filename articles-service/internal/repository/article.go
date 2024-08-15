package repository

import (
	"articles-service/internal/models"
	"context"
	"strconv"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ArticleRepo struct {
	ctx  context.Context
	conn *pgxpool.Pool
}

var _ ArticleRepoInterface = (*ArticleRepo)(nil)

func NewArticleRepo(ctx context.Context, conn *pgxpool.Pool) *ArticleRepo {
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

func (s *ArticleRepo) GetAll(dateSort string, page int, userId int,
) ([]*models.Article, error) {
	var orderQuery string
	if dateSort != "" {
		orderQuery += " ORDER BY timestamp " + dateSort
	}

	var whereQuery = ""
	if userId > 0 {
		whereQuery = " WHERE user_id = " + strconv.Itoa(userId)
	}

	const articlesPerPage = 10
	var skipArticles = articlesPerPage * (page)

	rows, err := s.conn.Query(s.ctx,
		`SELECT *
		FROM articles
		`+whereQuery+orderQuery+` LIMIT $1 OFFSET $2;`,
		articlesPerPage, skipArticles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articlesArr := make([]*models.Article, 0, articlesPerPage)
	err = pgxscan.ScanAll(&articlesArr, rows)

	return articlesArr, err
}
