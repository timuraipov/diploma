package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/internal/storage/db"
)

type OrderRepository struct {
	database db.DB
}

func NewOrderRepository(db db.DB) *OrderRepository {
	return &OrderRepository{
		database: db,
	}
}

func (o *OrderRepository) Save(ctx context.Context, order domain.Order) error {
	argsCheck := pgx.NamedArgs{
		"orderId": order.ID,
	}

	argsInsert := pgx.NamedArgs{
		"id":         order.ID,
		"status":     order.Status,
		"accrual":    order.Accrual,
		"userID":     order.UserID,
		"uploadedAt": order.UploadedAt,
	}

	return withTransaction(ctx, o.database.Pool, func(tx pgx.Tx) error {
		var orderFound domain.Order
		err := tx.QueryRow(ctx, SaveGetOrderByID, argsCheck).Scan(
			&orderFound.ID,
			&orderFound.Status,
			&orderFound.Accrual,
			&orderFound.UserID,
			&orderFound.UploadedAt,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				// Заказа нет, вставляем новый
				err = tx.QueryRow(ctx, SaveInsertNewOrder, argsInsert).Scan(&order.ID)
				if err != nil {
					return err
				}
				return nil
			}
			return err // другая ошибка
		}

		// Заказ уже существует
		if orderFound.ID == order.ID {
			if orderFound.UserID == order.UserID {
				return domain.ErrOrderAlreadyInProcessing
			}
			return domain.ErrOrderAlreadyProcessedByAnotherUser
		}

		return nil
	})
}

func (o *OrderRepository) GetAll(ctx context.Context, userID int64) ([]domain.Order, error) {

	args := pgx.NamedArgs{
		"userID": userID,
	}
	rows, err := o.database.Pool.Query(ctx, GetAllSelectAllOrderByUserID, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		err = rows.Scan(&order.ID, &order.Status, &order.Accrual, &order.UserID, &order.UploadedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *OrderRepository) GetUnhandledOrders(ctx context.Context) ([]domain.Order, error) {
	stmt := fmt.Sprintf(GetUnhandledOrdersGetOrder, domain.REGISTERED, domain.PROCESSING)
	rows, err := o.database.Pool.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		err = rows.Scan(&order.ID, &order.Status, &order.Accrual, &order.UserID, &order.UploadedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *OrderRepository) UpdateOrder(ctx context.Context, order domain.Order) error {

	args := pgx.NamedArgs{
		"id":      order.ID,
		"status":  order.Status,
		"accrual": order.Accrual,
	}

	return withTransaction(ctx, o.database.Pool, func(tx pgx.Tx) error {
		ct, err := tx.Exec(ctx, UpdateOrderUpdateOrderByID, args)
		if err != nil {
			return err
		}
		if ct.RowsAffected() == 0 {
			return errors.New("no rows updated")
		}
		return nil
	})
}
