package storage

import (
	"auth-service/internal/models"
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type AuthStorage struct {
	conn *pgx.Conn
}

var _ AuthStorageInterface = (*AuthStorage)(nil)

func (a *AuthStorage) GetApp(ctx context.Context, appID int32) (*models.App, error) {
	app := &models.App{}
	err := pgxscan.Get(ctx, a.conn, app, `
		SELECT *
		FROM apps 
		WHERE id=$1;`, appID)

	return app, err
}

func (a *AuthStorage) GetUser(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}

	err := pgxscan.Get(
		ctx, a.conn, user, `
		SELECT *
		FROM users 
		WHERE email=$1;`, email,
	)

	return user, err
}

func (a *AuthStorage) SaveUser(ctx context.Context, email string, passHash []byte,
) (int64, error) {
	var id int64
	err := a.conn.QueryRow(ctx,
		`INSERT INTO users(email, pass_hash) 
		VALUES($1, $2)
		RETURNING id;`,
		email, passHash,
	).Scan(&id)

	return id, err
}
