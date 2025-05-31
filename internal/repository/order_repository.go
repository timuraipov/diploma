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

	argsInsert := pgx.NamedArgs{
		"id":         order.ID,
		"status":     order.Status,
		"accrual":    order.Accrual,
		"userId":     order.UserId,
		"uploadedAt": order.UploadedAt,
	}

	return withTransaction(ctx, o.database.Pool, func(tx pgx.Tx) error {
		var orderFound domain.Order
		err := tx.QueryRow(ctx, stmtCheck, argsCheck).Scan(
			&orderFound.ID,
			&orderFound.Status,
			&orderFound.Accrual,
			&orderFound.UserId,
			&orderFound.UploadedAt,
		)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				// Заказа нет, вставляем новый
				err = tx.QueryRow(ctx, stmtExec, argsInsert).Scan(&order.ID)
				if err != nil {
					return err
				}
				return nil
			}
			return err // другая ошибка
		}

		// Заказ уже существует
		if orderFound.ID == order.ID {
			if orderFound.UserId == order.UserId {
				return domain.ErrOrderAlreadyInProcessing
			}
			return domain.ErrOrderAlreadyProcessedByAnotherUser
		}

		return nil
	})
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
	return orders, nil
}
func (o *orderRepository) GetUnhandledOrders(ctx context.Context) ([]domain.Order, error) {
	// const stmt = `SELECT id, status, accrual, user_id, uploaded_at FROM "order" WHERE status = 'REGISTERED' or status = 'PROCESSING'`
	var stmt = fmt.Sprintf(`SELECT id, status, accrual, user_id, uploaded_at FROM "order" WHERE status = '%s'  or status = '%s'`, domain.REGISTERED, domain.PROCESSING)
	rows, err := o.database.Pool.Query(ctx, stmt)
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
	return orders, nil
}
func (o *orderRepository) UpdateOrder(ctx context.Context, order domain.Order) error {
	const stmt = `UPDATE "order" SET status = @status, accrual = @accrual WHERE id = @id`
	args := pgx.NamedArgs{
		"id":      order.ID,
		"status":  order.Status,
		"accrual": order.Accrual,
	}

	return withTransaction(ctx, o.database.Pool, func(tx pgx.Tx) error {
		ct, err := tx.Exec(ctx, stmt, args)
		if err != nil {
			return err
		}
		if ct.RowsAffected() == 0 {
			return errors.New("no rows updated")
		}
		return nil
	})
}

// func (o *orderRepository) UpdateOrder(ctx context.Context, order domain.Order) error { // check is already processed by another user
// 	const stmt = `UPDATE "order" SET status = @status, accrual = @accrual WHERE id = @id`
// 	args := pgx.NamedArgs{
// 		"id":      order.ID,
// 		"status":  order.Status,
// 		"accrual": order.Accrual,
// 	}
// 	tx, err := o.database.Pool.BeginTx(ctx, pgx.TxOptions{})
// 	defer func() {
// 		if err != nil {
// 			tx.Rollback(ctx)
// 		} else {
// 			tx.Commit(ctx)
// 		}
// 	}()
// 	if err != nil {
// 		return err
// 	}
// 	row := tx.QueryRow(ctx, stmt, args)
// 	err = row.Scan()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

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
