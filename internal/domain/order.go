package domain

import (
	"context"
	"errors"
	"time"
)

const (
	REGISTERED string = "REGISTERED"
	PROCESSING        = "PROCESSING"
	INVALID           = "INVALID"
	PROCESSED         = "PROCESSED"
)

var (
	OrderAlreadyProcessedByAnotherUser = errors.New("Order with orderID already exists")
	OrderAlreadyInProcessing           = errors.New("Order with orderID already  processed")
)

type Order struct {
	ID         string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    int64     `json:"accrual,omitempty"`
	UserId     int64     `json:"-"`
	UploadedAt time.Time `json:"uploaded_at"`
}
type OrderRepository interface {
	Save(ctx context.Context, order Order) error
	GetAll(ctx context.Context, userId int64) ([]Order, error)
}

type OrderUseCase interface {
	Save(ctx context.Context, order Order) error
	GetAll(ctx context.Context, userId int64) ([]Order, error)
}
