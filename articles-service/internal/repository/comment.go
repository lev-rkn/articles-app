package repository

import (
	"articles-service/internal/lib/types"
	"articles-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type CommentRepo struct {
	ctx context.Context
	pg  *pgxpool.Pool
	rdb *redis.Client
}

var _ CommentRepoInterface = (*CommentRepo)(nil)

func NewCommentRepo(ctx context.Context, pg *pgxpool.Pool, rdb *redis.Client) *CommentRepo {
	return &CommentRepo{
		ctx: ctx,
		pg:  pg,
		rdb: rdb,
	}
}

func (r *CommentRepo) Create(comment *models.Comment) (int, error) {
	err := r.pg.QueryRow(r.ctx,
		`INSERT INTO comments (text, article_id, user_id) 
		VALUES ($1, $2, $3) 
		RETURNING id, timestamp;`,
		comment.Text,
		comment.ArticleId,
		comment.UserId).Scan(&comment.Id, &comment.Timestamp)

	articleIdStr := strconv.Itoa(comment.ArticleId)
	redisRes := r.rdb.Get(r.ctx, "comments:"+articleIdStr)
	// Если комментарии этой статьи закэшированы, добавляем этот комментарий туда же
	if redisRes.Err() == nil {
		var commentsArr []*models.Comment
		// берем комментарии из кэша
		result, _ := redisRes.Result()
		if err := json.Unmarshal([]byte(result), &commentsArr); err != nil {
			return -1, fmt.Errorf("unmarshal comments from cache")
		}

		commentsArr = append(commentsArr, comment)

		marshalled, err := json.Marshal(commentsArr)
		if err != nil {
			return -1, fmt.Errorf("comments marshalling: %w", err)
		}
		err = r.rdb.Set(r.ctx, "comments:"+articleIdStr, marshalled, 0).Err()
		if err != nil {
			return -1, fmt.Errorf("save comments to cache: %w", err)
		}
	}

	return comment.Id, err
}

func (r *CommentRepo) GetCommentsOnArticle(articleId int) ([]*models.Comment, error) {
	var commentsArr []*models.Comment
	articleIdStr := strconv.Itoa(articleId)

	redisRes := r.rdb.Get(r.ctx, "comments:"+articleIdStr)
	if redisRes.Err() != nil {
		// при отстутствии комментариев в кэше - берем из бд
		rows, err := r.pg.Query(r.ctx,
			`SELECT *
			FROM comments
			WHERE article_id=$1`, articleId)
		if err != nil {
			return nil, fmt.Errorf("get comments from db: %w", err)
		}

		if err := pgxscan.ScanAll(&commentsArr, rows); err != nil {
			return nil, fmt.Errorf("scan comment rows: %w", err)
		}
		rows.Close()

		if rows.CommandTag().RowsAffected() == 0 {
			return nil, types.ErrNoComments
		}

		// кэширование комментариев статьи
		marshalled, err := json.Marshal(commentsArr)
		if err != nil {
			return nil, fmt.Errorf("comments marshalling: %w", err)
		}
		err = r.rdb.Set(r.ctx, "comments:"+articleIdStr, marshalled, 0).Err()
		if err != nil {
			return nil, fmt.Errorf("save comments to cache: %w", err)
		}
	} else {
		// берем комментарии из кэша
		result, _ := redisRes.Result()
		if err := json.Unmarshal([]byte(result), &commentsArr); err != nil {
			return nil, fmt.Errorf("unmarshal comments from cache")
		}

		slog.Debug("get comments from cache")
	}

	return commentsArr, nil
}
