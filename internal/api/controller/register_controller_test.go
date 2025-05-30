package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/internal/repository/mocks"
	"github.com/timuraipov/diploma/internal/usecase"
	"github.com/timuraipov/diploma/pkg/logging"
	"go.uber.org/zap"
)

func setupTestRegisterRepository(t *testing.T) domain.UserRepository {
	userRepositoryMock := mocks.NewMockUserRepository()
	return &userRepositoryMock
}
func setupTestRegisterEnvironment(t *testing.T) (*RegisterController, func()) {
	logger, err := logging.NewZapLogger(zap.DebugLevel)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	cfg, err := bootstrap.MustLoad()
	assert.NoError(t, err)
	// Настройка тестовой базы данных
	userRepo := setupTestRegisterRepository(t)
	timeout := time.Duration(3000) * time.Second
	userUseCase := usecase.NewRegisterUsecase(logger, userRepo, timeout)

	// Создание контроллера
	registerController := NewRegisterController(logger, userUseCase, cfg)

	return registerController, func() {}
}

func TestRegisterController_RegisterUser(t *testing.T) {
	// Настройка окружения
	registerController, cleanup := setupTestRegisterEnvironment(t)
	defer cleanup()

	// Создаем тестовый HTTP-запрос
	user := domain.User{
		Login:    "testuser",
		Password: "testpassword",
	}
	body, err := json.Marshal(user)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Создаем тестовый HTTP-ответ
	rec := httptest.NewRecorder()

	// Вызываем метод контроллера
	registerController.Register(rec, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, rec.Code)

	// Проверяем, что пользователь был сохранен в базе данных
	savedUser, err := registerController.registerUsecase.GetByLogin(context.Background(), user.Login)
	require.NoError(t, err)
	assert.Equal(t, user.Login, savedUser.Login)
}
