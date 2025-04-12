package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/internal/storage/db"
)

type orderRepository struct {
	database db.DB
}

func NewOrderRepository(db db.DB) *orderRepository {
	return &orderRepository{
		database: db,
	}
}
func (o *orderRepository) Save(ctx context.Context, order domain.Order) error {
	const stmt = `INSERT INTO "order" (number, status,accrual,user_id,uploaded_at) VALUES (@number, @status,@accrual,@userId, @uploadedAt) returning (id)`
	args := pgx.NamedArgs{
		"number":     order.Number,
		"status":     order.Status,
		"accrual":    order.Accrual,
		"userId":     order.UserId,
		"uploadedAt": order.UploadedAt,
	}
	err := o.database.Pool.QueryRow(ctx, stmt, args).Scan(&order.Number)
	return err
}
func (o *orderRepository) GetAll(ctx context.Context, userId int64) ([]domain.Order, error) {
	const stmt = `SELECT number, status, accrual, user_id, uploaded_at FROM "order" WHERE user_id = @userId`
	args := pgx.NamedArgs{
		"userId": userId,
	}
	rows, err := o.database.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UserId, &order.UploadedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
