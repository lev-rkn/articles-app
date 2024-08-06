package middlewares

import (
	"articles-service/internal/config"
	"articles-service/internal/lib/types"
	"articles-service/internal/lib/utils"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// middleware, который проверяет jwt токен, лежащий в заголовке запроса, на валидность.
// в токене лежит закодированный аутентифицированный пользователь
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем JWT-токен из запроса
		tokenStr := extractBearerToken(c)
		if tokenStr == "" {
			// где-то и не нужна аутентификация
			return
		}

		// Парсим и валидируем токен, используя СЕКРЕТНЫЙ ключ
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Cfg.AuthGPRC.SecretKey), nil
		})
		if err != nil {
			utils.ErrorLog("failed to parse token", err)
			c.Set(types.KeyError, types.ErrInvalidToken)
			return
		}

		slog.Info("user authorized", slog.Any("claims", token))

		// Полученны данные сохраняем в контекст,
		// откуда его смогут получить следующие хэндлеры.
		c.Set(types.KeyUser, token)
	}
}

// Вынимает токен из заголовка запроса
func extractBearerToken(c *gin.Context) string {
	authHeader := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return ""
	}

	return splitToken[1]
}
