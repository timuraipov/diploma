package domain

import (
	"errors"
	"time"
)

var (
	ErrWithdrawAlreadyUsed = errors.New("order already used")
)

type Withdraw struct {
	ID          string    `json:"order" omitempty:"true"`
	Sum         float64   `json:"sum"`
	UserId      int64     `json:"-"`
	ProcessedAt time.Time `json:"processed_at"`
}
type WithdrawResponse struct {
	ID          string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
