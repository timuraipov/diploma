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
	client_mocks "github.com/timuraipov/diploma/internal/client/mocks"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/internal/repository/mocks"
	"github.com/timuraipov/diploma/internal/usecase"
	"github.com/timuraipov/diploma/pkg/logging"
	"go.uber.org/zap"
)

func setupOrderTestEnvironment(t *testing.T) (*OrderController, func()) {
	logger, err := logging.NewZapLogger(zap.DebugLevel)
	require.NoError(t, err)

	cfg := bootstrap.MustLoad()
	timeout := time.Duration(3000) * time.Second
	// Настройка mock-репозитория
	orderRepo := mocks.NewMockOrderRepository()
	balanceRepo := mocks.NewMockBalanceRepository()
	balanceUseCase := usecase.NewBalanceUseCase(logger, &balanceRepo, timeout)
	// Создание usecase

	client := client_mocks.NewMockClient(cfg.AccrualAddress)
	orderUseCase := usecase.NewOrderUseCase(logger, &orderRepo, client, balanceUseCase, timeout)

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
	body := []byte(string("39476215825"))
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
	// Создаем тестовый HTTP-ответ
	rec := httptest.NewRecorder()
	ctx := req.Context()
	ctx = context.WithValue(ctx, domain.UserIDHeader, int64(0))
	req = req.WithContext(ctx)

	// Вызываем метод контроллера
	orderController.Save(rec, req)
	// Создание запроса
	// Проверяем результат
	assert.Equal(t, http.StatusAccepted, rec.Code)
}

func TestOrderController_GetOrders(t *testing.T) {
	// Настройка окружения
	orderController, cleanup := setupOrderTestEnvironment(t)
	defer cleanup()
	ctx := context.Background()
	orders := []domain.Order{
		{
			ID:         "60412987384",
			Status:     domain.REGISTERED,
			Accrual:    0,
			UserID:     1,
			UploadedAt: time.Now(),
		},
		{
			ID:         "6011000990139424",
			Status:     domain.PROCESSING,
			Accrual:    0,
			UserID:     1,
			UploadedAt: time.Now(),
		},
		{
			ID:         "79927398713",
			Status:     domain.REGISTERED,
			Accrual:    0,
			UserID:     2,
			UploadedAt: time.Now(),
		},
	}
	for _, order := range orders {
		err := orderController.orderUseCase.Save(ctx, order)
		assert.NoError(t, err)
	}
	// Настройка mock-репозитория
	body := []byte(string("39476215825"))
	req := httptest.NewRequest(http.MethodGet, "/orders", bytes.NewReader(body))
	// Создаем тестовый HTTP-ответ
	rec := httptest.NewRecorder()
	ctx = req.Context()
	ctx = context.WithValue(ctx, domain.UserIDHeader, int64(1))
	req = req.WithContext(ctx)

	// Вызываем метод контроллера
	orderController.GetOrders(rec, req)
	// Создание запроса
	// Проверяем результат
	assert.Equal(t, http.StatusOK, rec.Code)
	var response []domain.Order
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(response))
}
