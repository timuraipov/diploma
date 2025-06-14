package domain

import (
	"context"
	"errors"
	"time"
)

const (
	REGISTERED string = "NEW"
	PROCESSING string = "PROCESSING"
	INVALID    string = "INVALID"
	PROCESSED  string = "PROCESSED"
)

var (
	ErrOrderAlreadyProcessedByAnotherUser = errors.New("order with orderID already exists")
	ErrOrderAlreadyInProcessing           = errors.New("order with orderID already  processed")
	ErrOrderNotFound                      = errors.New("order with orderID not found")
	ErrIncorrectOrderIDFormat             = errors.New("incorrect orderID format, should be 16 digits")
)

type Order struct {
	ID         string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UserID     int64     `json:"-"`
	UploadedAt time.Time `json:"uploaded_at"`
}
type OrderRepository interface {
	Save(ctx context.Context, order Order) error
	GetAll(ctx context.Context, UserID int64) ([]Order, error)
	GetUnhandledOrders(ctx context.Context) ([]Order, error)
	UpdateOrder(ctx context.Context, order Order) error
}

type OrderUseCase interface {
	Save(ctx context.Context, order Order) error
	GetAll(ctx context.Context, UserID int64) ([]Order, error)
	Accrual(ctx context.Context, order Order) (int, error)
	GetUnhandledOrders(ctx context.Context) ([]Order, error)
}
