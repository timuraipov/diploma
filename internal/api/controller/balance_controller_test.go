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

func setupBalanceTestEnvironment(t *testing.T) (*BalanceController, func()) {
	logger, err := logging.NewZapLogger(zap.DebugLevel)
	require.NoError(t, err)

	cfg, err := bootstrap.MustLoad()
	assert.NoError(t, err)
	// Настройка mock-репозитория
	balanceRepo := mocks.NewMockBalanceRepository()

	// Создание usecase
	timeout := time.Duration(3000) * time.Second
	balanceUseCase := usecase.NewBalanceUseCase(logger, &balanceRepo, timeout)

	// Создание контроллера
	balanceController := NewBalanceController(logger, balanceUseCase, cfg)

	require.NoError(t, err)
	require.NoError(t, err)

	return balanceController, func() {}
}

func TestBalanceController_GetBalance(t *testing.T) {
	balanceController, cleanup := setupBalanceTestEnvironment(t)
	defer cleanup()
	err := balanceController.balanceUseCase.UpdateBalance(context.Background(), int64(1), 100.1)
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodGet, "/balance", bytes.NewReader([]byte("")))
	// Создаем тестовый HTTP-ответ
	rec := httptest.NewRecorder()
	ctx := req.Context()
	ctx = context.WithValue(ctx, domain.UserIDHeader, int64(1))
	req = req.WithContext(ctx)

	// Вызываем метод контроллера
	balanceController.GetBalance(rec, req)
	// Создание запроса
	// Проверяем результат
	assert.Equal(t, http.StatusOK, rec.Code)
	var response domain.BalanceResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 100.1, response.Current)
}

func TestBalanceController_Withdraw(t *testing.T) {
	balanceController, cleanup := setupBalanceTestEnvironment(t)
	defer cleanup()
	err := balanceController.balanceUseCase.UpdateBalance(context.Background(), int64(1), 101)
	assert.NoError(t, err)
	withdrawRequest := domain.WithdrawRequest{
		Order: "neworder",
		Sum:   100.12,
	}
	body, err := json.Marshal(withdrawRequest)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewReader(body))
	// Создаем тестовый HTTP-ответ
	rec := httptest.NewRecorder()
	ctx := req.Context()
	ctx = context.WithValue(ctx, domain.UserIDHeader, int64(1))
	req = req.WithContext(ctx)

	// Вызываем метод контроллера
	balanceController.Withdraw(rec, req)
	// Создание запроса
	// Проверяем результат
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestBalanceController_Withdrawals(t *testing.T) {
	balanceController, cleanup := setupBalanceTestEnvironment(t)
	defer cleanup()
	err := balanceController.balanceUseCase.UpdateBalance(context.Background(), int64(1), 101)
	assert.NoError(t, err)
	withdrawRequests := []domain.Withdraw{
		{
			OrderID:     "order1",
			Sum:         20.1,
			UserID:      1,
			ProcessedAt: time.Now(), // time.Now().Format(time.RFC3339),

		},
		{
			OrderID:     "order2",
			Sum:         11.12,
			UserID:      1,
			ProcessedAt: time.Now(),
		},
	}
	for _, testCase := range withdrawRequests {
		err = balanceController.balanceUseCase.Withdraw(context.Background(), testCase)
		assert.NoError(t, err)
	}

	req := httptest.NewRequest(http.MethodGet, "/withdrawals", bytes.NewReader([]byte("")))
	// Создаем тестовый HTTP-ответ
	rec := httptest.NewRecorder()
	ctx := req.Context()
	ctx = context.WithValue(ctx, domain.UserIDHeader, int64(1))
	req = req.WithContext(ctx)

	// Вызываем метод контроллера
	balanceController.Withdrawals(rec, req)
	// Создание запроса
	// Проверяем результат
	assert.Equal(t, http.StatusOK, rec.Code)
	var responseWithdrawals []domain.WithdrawResponse
	err = json.Unmarshal(rec.Body.Bytes(), &responseWithdrawals)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(responseWithdrawals))
}
