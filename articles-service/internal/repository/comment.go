package repository

import (
	"articles-service/internal/models"
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type CommentRepo struct {
	ctx  context.Context
	conn *pgx.Conn
}

var _ CommentRepoInterface = (*CommentRepo)(nil)

func NewCommentRepo(ctx context.Context, conn *pgx.Conn) *CommentRepo {
	return &CommentRepo{
		ctx:  ctx,
		conn: conn,
	}
}

func (r *CommentRepo) Create(comment *models.Comment) (int, error) {
	var id int
	err := r.conn.QueryRow(r.ctx,
		`INSERT INTO comments (text, article_id, user_id) 
		VALUES ($1, $2, $3) 
		RETURNING id;`,
		comment.Text,
		comment.ArticleId,
		comment.UserId).Scan(&id)

	return id, err
}

func (r *CommentRepo) GetCommentsOnArticle(articleId int) ([]*models.Comment, error) {
	rows, err := r.conn.Query(r.ctx,
		`SELECT *
		FROM comments
		WHERE article_id=$1`, articleId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var commentsArr []*models.Comment
	err = pgxscan.ScanAll(&commentsArr, rows)

	return commentsArr, err
}
