package storage

import (
	"auth-service/internal/lib/types"
	"auth-service/internal/models"
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, types.ErrAppNotFound
		}

		return nil, err
	}

	return app, nil
}

func (a *AuthStorage) GetUser(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}

	err := pgxscan.Get(
		ctx, a.conn, user, `
		SELECT *
		FROM users 
		WHERE email=$1;`, email,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, types.ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
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

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return -1, types.ErrUserExists
			}
		}

		return -1, err
	}

	return id, nil
}

func (a *AuthStorage) SaveRefreshSession(
	ctx context.Context,
	refresh_token string,
	fingerprint string,
	userId int,
) error {
	_, err := a.conn.Exec(ctx,
		`INSERT INTO refreshSessions(refresh_token, fingerprint, user_id) 
		VALUES ($1, $2, $3);`,
		refresh_token, fingerprint, userId,
	)

	return err
}

func (a *AuthStorage) GetRefreshSession(
	ctx context.Context,
	refreshToken string,
) (*models.RefreshSession, error) {
	session := &models.RefreshSession{}

	err := a.conn.QueryRow(
		ctx, `
		SELECT refresh_token, fingerprint, user_email
		FROM refreshSessions 
		WHERE refresh_token=$1;`,
		refreshToken,
	).Scan(&session.RefreshToken, &session.Fingerprint, &session.UserEmail)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, types.ErrRefreshTokenNotValid
		}

		return nil, err
	}

	return session, nil
}
