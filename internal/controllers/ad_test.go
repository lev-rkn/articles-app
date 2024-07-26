package controllers

import (
	"ads-service/internal/models"
	"ads-service/internal/service/mocks"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAd(t *testing.T) {
	// добавить инициализацию логгера, чтобы логировать в файл
	testAd := &models.Ad{
		Title:       "title",
		Description: "description",
		Price:       100,
		Photos:      []string{"photo1", "photo2"},
	}

	tests := []struct {
		name         string
		mockExpect   func(adService *mocks.AdServiceInterface)
		ad           *models.Ad
		expectOutput string
		code         int
	}{
		{
			name: "Валидный кейс",
			mockExpect: func(adService *mocks.AdServiceInterface) {
				adService.On("Create", testAd).Return(3, nil)
			},
			ad:           testAd,
			expectOutput: `{"id": 3}`,
			code:         http.StatusCreated,
		},
		{
			name:       "Заголовок длиной более 200 символов",
			mockExpect: func(adService *mocks.AdServiceInterface) {},
			ad: &models.Ad{
				Photos:      testAd.Photos,
				Title:       strings.Repeat("a", 201),
				Description: testAd.Description,
				Price:       testAd.Price,
			},
			expectOutput: ErrInvalidTitle.Error() + "\n",
			code:         http.StatusBadRequest,
		},
		{
			name:       "Заголовок пуст",
			mockExpect: func(adService *mocks.AdServiceInterface) {},
			ad: &models.Ad{
				Photos:      testAd.Photos,
				Title:       "",
				Description: testAd.Description,
				Price:       testAd.Price,
			},
			expectOutput: ErrEmptyTitle.Error() + "\n",
			code:         http.StatusBadRequest,
		},
		{
			name:       "Длина описания больше 1000 символов",
			mockExpect: func(adService *mocks.AdServiceInterface) {},
			ad: &models.Ad{
				Photos:      testAd.Photos,
				Title:       testAd.Title,
				Description: strings.Repeat("a", 1001),
				Price:       testAd.Price,
			},
			expectOutput: ErrInvalidDescription.Error() + "\n",
			code:         http.StatusBadRequest,
		},
		{
			name:       "Попытка загрузить более чем 3 ссылки на фото",
			mockExpect: func(adService *mocks.AdServiceInterface) {},
			ad: &models.Ad{
				Photos:      []string{"photo1", "photo2", "photo3", "photo4"},
				Title:       testAd.Title,
				Description: testAd.Description,
				Price:       testAd.Price,
			},
			expectOutput: ErrInvalidPhotos.Error() + "\n",
			code:         http.StatusBadRequest,
		},
		{
			name: "Случайная ошибка от сервиса",
			mockExpect: func(adService *mocks.AdServiceInterface) {
				adService.On("Create", testAd).Return(0, errors.New("some error"))
			},
			ad:           testAd,
			expectOutput: errors.New("some error").Error() + "\n",
			code:         http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		// инициализация контекста
		ctx := context.Background()
		// инициализация мок-сервиса
		mockService := &mocks.AdServiceInterface{}
		// инициализация тестируемого контроллера
		mux := http.NewServeMux()
		adController := InitAdController(ctx, mockService, mux)

		// превращаем тестовое объявление в JSON
		marshalled, _ := json.Marshal(test.ad)
		// создаем тестовый запрос, который будем пихать в тестируемый контроллер
		req, _ := http.NewRequest("POST", "/ad/create/", strings.NewReader(string(marshalled)))
		// создаем подставной обработчик, который будет слушать ответ тестируемого контроллера
		w := httptest.NewRecorder()
		// саем моку ожидаемое поведение
		test.mockExpect(mockService)
		// запускаем тестируемый контроллер, который должен записать свой ответ в
		// инициализированный выше подставной обработчик
		adController.Create(w, req)

		// проверка соответствия кода ответа
		assert.Equal(t, test.code, w.Code)
		// проверка записи котроллера
		assert.Equal(t, test.expectOutput, w.Body.String())
		// проверка, что ожадаемое моком поведение в точности выполнено
		mockService.AssertExpectations(t)
	}
}
