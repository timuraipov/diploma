package domain

import (
	"context"
	"errors"
	"time"
)

const (
	NEW        string = "new"
	PROCESSING        = "processing"
	INVALID           = "invalid"
	PROCESSED         = "processed"
)

var (
	OrderAlreadyProcessedByAnotherUser = errors.New("Order with orderID already exists")
	OrderAlreadyInProcessing           = errors.New("Order with orderID already  processed")
)

type Order struct {
	ID         string    `json:"id"`
	Status     string    `json:"status"`
	Accrual    int64     `json:"accrual"`
	UserId     int64     `json:"user_id" ommitepmty:"true"`
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
