package controller

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/repository/mocks"
	"github.com/timuraipov/diploma/internal/usecase"
	"github.com/timuraipov/diploma/pkg/logging"
	"go.uber.org/zap"
)

func setupOrderTestEnvironment(t *testing.T) (*OrderController, func()) {
	logger, err := logging.NewZapLogger(zap.DebugLevel)
	require.NoError(t, err)

	cfg, err := bootstrap.MustLoad()
	require.NoError(t, err)

	// Настройка mock-репозитория
	orderRepo := mocks.NewMockOrderRepository()

	// Создание usecase
	timeout := time.Duration(3000) * time.Second
	orderUseCase := usecase.NewOrderUseCase(logger, &orderRepo, timeout)

	// Создание контроллера
	orderController := NewOrderController(logger, orderUseCase, cfg)

	require.NoError(t, err)
	require.NoError(t, err)

	return orderController, func() {}
}
func TestOrderController_CreateOrder(t *testing.T) {
	// Настройка окружения
	orderController, cleanup := setupOrderTestEnvironment(t)
	defer cleanup()
	// Настройка mock-репозитория
	body := []byte(string("orderID"))
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
	// Создаем тестовый HTTP-ответ
	rec := httptest.NewRecorder()
	ctx := req.Context()
	ctx = context.WithValue(ctx, "x-user-id", int64(0))
	req = req.WithContext(ctx)

	// Вызываем метод контроллера
	orderController.Save(rec, req)
	// Создание запроса
	// Проверяем результат
	assert.Equal(t, http.StatusAccepted, rec.Code)

}
