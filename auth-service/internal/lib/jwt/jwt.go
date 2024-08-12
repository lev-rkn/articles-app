package jwt

import (
	"auth-service/internal/config"
	"auth-service/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewTokenPair(user *models.User, app *models.App) (*models.TokenPair, error) {
	accessToken, err := generateAccessToken(user, app, config.Cfg.AccessTokenTTL)
	if err != nil {
		return nil, err
	}
	refreshToken, err := generateRefreshToken(app, config.Cfg.RefreshTokenTTL)
	if err != nil {
		return nil, err
	}

	tokenPain := &models.TokenPair{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}

	return tokenPain, nil
}

func generateAccessToken(user *models.User, app *models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
    claims["uid"] = user.ID  
    claims["email"] = user.Email  
    claims["exp"] = time.Now().Add(duration).Unix()  
    claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func generateRefreshToken(app *models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
