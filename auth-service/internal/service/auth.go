package service

import (
	"auth-service/internal/lib/jwt"
	"auth-service/internal/lib/types"
	"auth-service/internal/models"
	"auth-service/internal/storage"
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authStorage storage.AuthStorageInterface
}

func (s *AuthService) RegisterNewUser(
	ctx context.Context, email string, pass string) (int64, error) {

	// Генерируем хэш и соль для пароля.
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return -1, fmt.Errorf("generate password hash: %w", err)
	}

	// Сохраняем пользователя в БД
	id, err := s.authStorage.SaveUser(ctx, email, passHash)
	if err != nil {
		return -1, fmt.Errorf("authStorage.SaveUser: %w", err)
	}

	return id, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
	appID int32,
	fingerprint string,
) (*models.TokenPair, error) {
	// Достаём пользователя из БД
	user, err := s.authStorage.GetUser(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("authStorage.GetUser: %w", err)
	}

	// Проверяем пароль пользователя на соответствие
	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
	if err != nil {
		return nil, types.ErrInvalidCredentials
	}

	// Получаем информацию о приложении
	app, err := s.authStorage.GetApp(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("authStorage.GetApp: %w", err)
	}

	// Генерируем токен доступа и рефреш токен
	tokenPair, err := jwt.NewTokenPair(user, app)
	if err != nil {
		return nil, fmt.Errorf("jwt.NewTokenPair: %w", err)
	}

	// Сохраняем рефреш токен в сессию пользователя
	session := &models.RefreshSession{
		RefreshToken: tokenPair.RefreshToken,
		Fingerprint:  fingerprint,
		UserEmail:    user.Email,
		AppId:        appID,
	}
	err = s.authStorage.SaveRefreshSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("authStorage.SaveRefreshSession: %w", err)
	}

	return tokenPair, nil
}

func (s *AuthService) RefreshToken(
	ctx context.Context,
	refreshToken string,
	fingerprint string,
) (*models.TokenPair, error) {
	session, err := s.authStorage.GetRefreshSession(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("authStorage.GetRefreshSession: %w", err)
	}
	if fingerprint != session.Fingerprint {
		return nil, types.ErrUnidentifiedDevice
	}

	// TODO: делаем 2 лишних запроса на получение юзера и данных о приложении
	// нужно это дело кэшировать, теоритически данные можно брать из токена.

	// Достаём пользователя из БД
	user, err := s.authStorage.GetUser(ctx, session.UserEmail)
	if err != nil {
		errors.Join()
		return nil, fmt.Errorf("authStorage.GetUser: %w", err)
	}
	// Получаем информацию о приложении
	app, err := s.authStorage.GetApp(ctx, session.AppId)
	if err != nil {
		return nil, fmt.Errorf("authStorage.GetApp: %w", err)
	}
	// Генерируем токен доступа и рефреш токен
	tokenPair, err := jwt.NewTokenPair(user, app)
	if err != nil {
		return nil, fmt.Errorf("jwt.NewTokenPair: %w", err)
	}
	// Сохраняем новую сессию
	session.RefreshToken = tokenPair.RefreshToken
	err = s.authStorage.SaveRefreshSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("authStorage.SaveRefreshSession: %w", err)
	}

	return tokenPair, nil
}
