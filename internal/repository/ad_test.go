package repository

import (
	"ads-service/internal/models"
	"ads-service/internal/repository/mocks"
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	testAd := &models.Ad{
		Id:          1,
		Title:       "title",
		Price:       100,
		Photos:      []string{"photo1", "photo2"},
		Description: "description",
	}
	testQ := "INSERT INTO advertisements (title, description, price, photos) VALUES ($1, $2, $3, $4) RETURNING id;"

	testCases := []struct {
		name       string
		mockExpect func(ctx context.Context, userRepo *mocks.PgConn)
		err        error
	}{
		{
			name: "Успешный кейс, валидный ввод",
			mockExpect: func(ctx context.Context, userRepo *mocks.PgConn) {
				userRepo.On("QueryRow", ctx, testQ,
					testAd.Title, testAd.Description, testAd.Price, testAd.Photos).Return(userRepo.).
					On("Scan", mock.Anything).Return(nil)
			},
			err: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// инициализация
			ctx := context.Background()
			mockConn := &mocks.PgConn{}
			adRepo := NewAdRepo(ctx, mockConn)
			// включаем ожидание мока
			testCase.mockExpect(ctx, mockConn)
			// вызов
			_, err := adRepo.Create(testAd)
			// проверка
			assert.Equal(t, testCase.err, err)
			mockConn.AssertExpectations(t)
		})

	}
}
