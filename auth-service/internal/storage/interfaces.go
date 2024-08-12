package storage

import (
	"auth-service/internal/models"
	"context"
)

type AuthStorageInterface interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
	GetUser(ctx context.Context, email string) (*models.User, error)
	GetApp(ctx context.Context, appID int32) (*models.App, error)
	SaveRefreshSession(
		ctx context.Context,
		refresh_token string,
		fingerprint string,
		userId int,
	) error
	GetRefreshSession(
		ctx context.Context,
		refreshToken string,
	) (*models.RefreshSession, error)
}
