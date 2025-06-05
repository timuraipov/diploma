package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/internal/repository/mocks"
	"github.com/timuraipov/diploma/internal/usecase"
	"github.com/timuraipov/diploma/pkg/logging"
	"go.uber.org/zap"
)

func setupLoginTestEnvironment(t *testing.T) (*LoginController, func()) {
	logger, err := logging.NewZapLogger(zap.DebugLevel)
	require.NoError(t, err)

	cfg := bootstrap.MustLoad()

	// Настройка mock-репозитория
	userRepo := mocks.NewMockUserRepository()

	// Создание usecase
	loginUseCase := usecase.NewLoginUsecase(logger, &userRepo, cfg)

	// Создание контроллера
	loginController := NewLoginController(logger, loginUseCase, cfg)
	userUseCase := usecase.NewRegisterUsecase(logger, &userRepo, cfg)
	userRequest := domain.RegisterRequest{
		Login:    "testuser",
		Password: "testpassword",
	}
	_, err = userUseCase.Create(context.Background(), userRequest)
	require.NoError(t, err)
	return loginController, func() {}
}

func TestLoginController_LoginUser(t *testing.T) {
	// Настройка окружения
	loginController, cleanup := setupLoginTestEnvironment(t)
	defer cleanup()

	// Создаем тестового пользователя'''
	mockUser := domain.User{
		ID:       1,
		Login:    "testuser",
		Password: "testpassword",
	}

	// Настраиваем mock-репозиторий
	mockRepo := mocks.NewMockUserRepository()
	err := mockRepo.Create(context.Background(), &mockUser)
	assert.NoError(t, err)
	// создаем регистрацию

	// Создаем тестовый HTTP-запрос
	loginRequest := map[string]string{
		"login":    mockUser.Login,
		"password": mockUser.Password,
	}
	body, err := json.Marshal(loginRequest)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Создаем тестовый HTTP-ответ
	rec := httptest.NewRecorder()

	// Вызываем метод контроллера
	loginController.Login(rec, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, rec.Code)

	// Проверяем, что в ответе есть токен (или другая ожидаемая информация)
	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.NotEmpty(t, response["accessToken"])
	assert.NotEmpty(t, response["refreshToken"])
}
