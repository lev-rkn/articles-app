package service

import (
	"auth-service/internal/models"
	"context"
)

type AuthServiceInterface interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int32,
		fingerprint string,
	) (tokenPair *models.TokenPair, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
	RefreshToken(
		ctx context.Context,
		refreshToken string,
		fingerprint string,
	) (*models.TokenPair, error)
}
