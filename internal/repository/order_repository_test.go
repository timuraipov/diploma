package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/internal/storage/db"
)

var testRepositoryOrder *orderRepository

func setupTestDBOrder(t *testing.T) {

	dsn := os.Getenv("DATABASE_URI")
	if dsn == "" {
		t.Fatal("DATABASE_URI is not set")
	}

	db, err := db.NewDB(context.Background(), dsn)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	_, err = db.Pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS "order"(
    id          VARCHAR(30) PRIMARY KEY ,
    status      VARCHAR(20) NOT NULL,
    user_id     INT NOT NULL,
    accrual     INT DEFAULT 0,
    uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now() 
);
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
	testRepositoryOrder = NewOrderRepository(*db)
}
func cleanupTestDBOrder(t *testing.T) {
	_, err := testRepositoryOrder.database.Pool.Exec(context.Background(), `DROP TABLE "order"`)
	if err != nil {
		t.Fatalf("failed to drop table: %v", err)
	}
}
func TestOrderRepository_SaveAndGetAll(t *testing.T) {
	// Подготовка тестовой базы данных
	setupTestDBOrder(t)
	defer cleanupTestDBOrder(t)

	// Создаем тестовый заказ
	order := domain.Order{
		ID:         "test-order-id",
		Status:     "new",
		Accrual:    100,
		UserId:     1,
		UploadedAt: time.Now(),
	}

	// Тестируем метод Save
	err := testRepositoryOrder.Save(context.Background(), order)
	require.NoError(t, err, "Save method failed")

	// Тестируем метод GetAll
	orders, err := testRepositoryOrder.GetAll(context.Background(), order.UserId)
	require.NoError(t, err, "GetAll method failed")
	require.Len(t, orders, 1, "Expected one order in the result")
	// Проверяем данные заказа
	savedOrder := orders[0]
	assert.Equal(t, order.ID, savedOrder.ID)
	assert.Equal(t, order.Status, savedOrder.Status)
	assert.Equal(t, order.Accrual, savedOrder.Accrual)
	assert.Equal(t, order.UserId, savedOrder.UserId)
	assert.WithinDuration(t, order.UploadedAt, savedOrder.UploadedAt, time.Second)

	err = testRepositoryOrder.Save(context.Background(), order)
	require.ErrorIs(t, err, domain.ErrOrderAlreadyInProcessing)
	order.UserId = 0
	err = testRepositoryOrder.Save(context.Background(), order)
	require.ErrorIs(t, err, domain.ErrOrderAlreadyProcessedByAnotherUser)
}
