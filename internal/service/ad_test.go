package services

import (
	"ads-service/internal/models"
	"ads-service/internal/repository"
	"ads-service/internal/repository/mocks"
	"ads-service/logger"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetOne(t *testing.T) {

	input := &models.Ad{
		Id:          1,
		Title:       "title",
		Price:       100,
		Photos:      []string{"photo1", "photo2"},
		Description: "description",
	}

	testCases := []struct {
		name         string
		expectations func(userRepo *mocks.AdRepo)
		input        *models.Ad
		err          error
	}{
		{
			name: "объявление найдено",
			expectations: func(userRepo *mocks.AdRepo) {
				userRepo.On("GetOne", input.Id).Return(input, nil)
			},
			input: input,
			err:   nil,
		},
		{
			name: "объявление не найдено",
			expectations: func(userRepo *mocks.AdRepo) {
				userRepo.On("GetOne", input.Id).Return(nil, pgx.ErrNoRows)
			},
			input: input,
			err:   errAdNotFound,
		},
		{name: "Получена любая другая ошибка",
			expectations: func(userRepo *mocks.AdRepo) {
				userRepo.On("GetOne", input.Id).Return(nil, errors.New("some error"))
			},
			input: input,
			err:   errors.New("some error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			logger := logger.New()
			mockAdRepo := &mocks.AdRepo{}
			repository := &repository.Repository{Ad: mockAdRepo}
			service := NewService(repository, logger)

			testCase.expectations(mockAdRepo)
			_, err := service.Ad.GetOne(input.Id)

			assert.Equal(t, testCase.err, err)
		})
	}
}
