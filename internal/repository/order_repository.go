package repository

import (
	"context"
	"errors"
	"fmt"

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
	const stmtCheck = `SELECT id, status, accrual, user_id, uploaded_at 
	FROM "order" 
	WHERE id = @orderId`
	argsCheck := pgx.NamedArgs{
		"orderId": order.ID,
	}
	const stmtExec = `INSERT INTO "order" 
	(id, status, accrual, user_id, uploaded_at) 
	VALUES (@id, @status, @accrual, @userId, @uploadedAt) 
	RETURNING id`

	args := pgx.NamedArgs{
		"id":         order.ID,
		"status":     order.Status,
		"accrual":    order.Accrual,
		"userId":     order.UserId,
		"uploadedAt": order.UploadedAt,
	}
	tx, err := o.database.Pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	if err != nil {
		return err
	}
	row := tx.QueryRow(ctx, stmtCheck, argsCheck)
	var orderFound domain.Order
	err = row.Scan(&orderFound.ID, &orderFound.Status, &orderFound.Accrual, &orderFound.UserId, &orderFound.UploadedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = o.database.Pool.QueryRow(ctx, stmtExec, args).Scan(&order.ID)
			if err != nil {
				return err
			}
		}
		return err
	}
	if orderFound.ID != "" && orderFound.ID == order.ID {
		if orderFound.UserId == order.UserId {
			return domain.OrderAlreadyInProcessing
		}
		return domain.OrderAlreadyProcessedByAnotherUser
	}

	return nil
}
func (o *orderRepository) GetAll(ctx context.Context, userId int64) ([]domain.Order, error) {
	const stmt = `SELECT id, status, accrual, user_id, uploaded_at FROM "order" WHERE user_id = @userId`
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
		err = rows.Scan(&order.ID, &order.Status, &order.Accrual, &order.UserId, &order.UploadedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	fmt.Println(orders)
	return orders, nil
}

// if err != nil {
// 	if pgErr, ok := err.(*pgconn.PgError); ok {
// 		if pgErr.Code == pgerrcode.UniqueViolation {
// 			return domain.OrderAlready
// 		} else {
// 			fmt.Printf("🎯 PgError: %s (Code: %s)\n", pgErr.Message, pgErr.Code)
// 		}
// 	} else {
// 		fmt.Printf("❗ Не PgError: %T - %v\n", err, err)
// 	}
// }
