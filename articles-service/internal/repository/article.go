package repository

import (
	"articles-service/internal/lib/types"
	"articles-service/internal/models"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type ArticleRepo struct {
	ctx context.Context
	pg  *pgxpool.Pool
	rdb *redis.Client
}

var _ ArticleRepoInterface = (*ArticleRepo)(nil)

func NewArticleRepo(ctx context.Context, pg *pgxpool.Pool, rdb *redis.Client) *ArticleRepo {
	return &ArticleRepo{
		ctx: ctx,
		pg:  pg,
		rdb: rdb,
	}
}

func (r *ArticleRepo) Create(article *models.Article) (int, error) {
	var id int
	err := r.pg.QueryRow(r.ctx,
		`INSERT INTO articles (title, description, photos, user_id) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id;`,
		article.Title,
		article.Description,
		article.Photos,
		article.UserId).Scan(&id)

	return id, err
}

func (r *ArticleRepo) GetOne(id int) (*models.Article, error) {
	idStr := strconv.Itoa(id)
	article := &models.Article{}
	// попытка получить статью из кэша
	redisRes := r.rdb.Get(r.ctx, "article:"+idStr)
	if redisRes.Err() == nil {
		if err := redisRes.Scan(article); err != nil {
			return nil, fmt.Errorf("redisRes.Scan(article): %w", err)
		}
		slog.Debug("get article from cache")

		return article, nil
	}

	// если статьи нету в кэше, то делаем запрос в бд
	err := pgxscan.Get(
		r.ctx, r.pg, article, `
			SELECT *
			FROM articles 
			WHERE id=$1;`, id,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, types.ErrArticleNotFound
		}
		return nil, fmt.Errorf("get article from pg: %w", err)
	}
	// кэшируем полученную из бд статью
	err = r.rdb.Set(r.ctx, "article:"+idStr, article, 0).Err()
	if err != nil {
		return nil, fmt.Errorf("save article to cache: %w", err)
	}

	return article, nil
}

func (r *ArticleRepo) GetAll(dateSort string, page int, userId int,
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

	rows, err := r.pg.Query(r.ctx,
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
