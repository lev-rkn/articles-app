package controllers

import (
	"ads-service/internal/config"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrFailedIsAdminCheck = errors.New("failed to check if user is admin")
)

// middleware, который проверяет jwt токен, лежащий в заголовке запроса, на валидность.
// в токене лежит закодированный авторизованный пользователь
func AuthMiddleware(
	handlerFunc func(w http.ResponseWriter, r *http.Request),
) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем JWT-токен из запроса
		tokenStr := extractBearerToken(r)

		// Парсим и валидируем токен, используя СЕКРЕТНЫЙ ключ
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Cfg.AuthGPRC.SecretKey), nil
		})
		if err != nil {
			slog.Warn("failed to parse token", "err", err.Error())

			ctx := context.WithValue(r.Context(), "error", ErrInvalidToken)
			handlerFunc(w, r.WithContext(ctx))

			return
		}

		slog.Info("user authorized", slog.Any("claims", token))

		// Полученны данные сохраняем в контекст,
		// откуда его смогут получить следующие хэндлеры.
		ctx := context.WithValue(r.Context(), "user", token)

		handlerFunc(w, r.WithContext(ctx))
	})
}

// Вынимает токен из заголовка запроса
func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return ""
	}

	return splitToken[1]
}
