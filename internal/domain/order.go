package domain

import (
	"context"
)

const (
	NEW        string = "new"
	PROCESSING        = "processing"
	INVALID           = "invalid"
	PROCESSED         = "processed"
)

type Order struct {
	Number     string `json:"number"`
	Status     string `json:"status"`
	Accrual    int64  `json:"accrual"`
	UserId     int64  `json:"user_id" ommitepmty:"true"`
	UploadedAt string `json:"uploaded_at"`
}
type OrderRepository interface {
	Save(ctx context.Context, order Order) error
	GetAll(ctx context.Context, userId int64) ([]Order, error)
}

type OrderUseCase interface {
	Save(ctx context.Context, order Order) error
	GetAll(ctx context.Context, userId int64) ([]Order, error)
}
