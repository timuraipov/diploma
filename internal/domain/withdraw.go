package domain

import "time"

type Withdraw struct {
	ID          string    `json:"order" ommitempty:"true"`
	Sum         float64   `json:"sum"`
	UserId      int64     `json:"user_id" ommitempty:"true"`
	ProcessedAt time.Time `json:"processed_at"`
}
type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
type WithdrawResponse struct {
	Withdrawals []Withdraw
}
