package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/internal/storage/db"
)

type balanceRepository struct {
	database db.DB
}

func NewBalanceRepository(db db.DB) *balanceRepository {
	return &balanceRepository{
		database: db,
	}
}

func (b *balanceRepository) GetBalance(ctx context.Context, userID int64) (domain.Balance, error) {
	var balance domain.Balance
	stmt := `SELECT id,balance,user_id,created_at,updated_at FROM "balance" where user_id=@userID`
	args := pgx.NamedArgs{
		"userID": userID,
	}
	err := b.database.Pool.QueryRow(ctx, stmt, args).Scan(&balance.ID, &balance.Balance, &balance.UserID, &balance.CreatedAt, &balance.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Balance{}, nil
		}
	}
	return balance, err
}

func (b *balanceRepository) UpdateBalance(ctx context.Context, userID int64, accrual float64) error {
	return withTransaction(ctx, b.database.Pool, func(tx pgx.Tx) error {
		const stmtSelect = `
		SELECT balance FROM "balance" WHERE user_id = @userID`
		selectArgs := pgx.NamedArgs{
			"userID": userID,
		}

		var currentBalance float64
		err := tx.QueryRow(ctx, stmtSelect, selectArgs).Scan(&currentBalance)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				// Баланс не существует, создаём новую запись
				const stmtInsert = `
				INSERT INTO "balance" (balance, user_id, created_at, updated_at)
				VALUES (@balance, @userID, @createdAt, @updatedAt)
				RETURNING id`
				now := time.Now()
				insertArgs := pgx.NamedArgs{
					"balance":   accrual,
					"userID":    userID,
					"createdAt": now,
					"updatedAt": now,
				}
				var id int64
				err = tx.QueryRow(ctx, stmtInsert, insertArgs).Scan(&id)
				return err
			}
			// Другая ошибка при выборке
			return err
		}

		// Баланс существует, обновляем его
		const stmtUpdate = `
		UPDATE "balance"
		SET balance = balance + @accrual, updated_at = @updatedAt
		WHERE user_id = @userID
		RETURNING balance`
		updateArgs := pgx.NamedArgs{
			"accrual":   accrual,
			"userID":    userID,
			"updatedAt": time.Now(),
		}
		var updatedBalance float64
		err = tx.QueryRow(ctx, stmtUpdate, updateArgs).Scan(&updatedBalance)
		return err
	})
}

func (b *balanceRepository) Withdraw(ctx context.Context, withdraw domain.Withdraw) error {
	return withTransaction(ctx, b.database.Pool, func(tx pgx.Tx) error {
		// Проверка: был ли уже использован этот ID
		const stmtWithdrawExist = `
		SELECT id FROM "withdraw" WHERE id = @orderID`
		withdrawExistsArgs := pgx.NamedArgs{
			"orderID": withdraw.ID,
		}
		var withdrawAlreadyExist string
		err := tx.QueryRow(ctx, stmtWithdrawExist, withdrawExistsArgs).Scan(&withdrawAlreadyExist)
		if err == nil {
			return domain.ErrWithdrawAlreadyUsed
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}

		// Получение текущего баланса
		const stmtBalanceSelect = `SELECT balance.balance FROM "balance" WHERE user_id = @userID`
		balanceArgs := pgx.NamedArgs{
			"userID": withdraw.UserID,
		}
		var balance float64
		err = tx.QueryRow(ctx, stmtBalanceSelect, balanceArgs).Scan(&balance)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		if errors.Is(err, pgx.ErrNoRows) || balance < withdraw.Sum {
			return domain.ErrInsufficientFunds
		}

		// Списание средств
		const stmtBalanceUpdate = `
		UPDATE "balance" 
		SET balance = balance - @accrual, updated_at = @updatedAt
		WHERE user_id = @userID
		RETURNING balance`
		balanceUpdateArgs := pgx.NamedArgs{
			"userID":    withdraw.UserID,
			"updatedAt": time.Now(),
			"accrual":   withdraw.Sum,
		}
		var finalBalance float64
		err = tx.QueryRow(ctx, stmtBalanceUpdate, balanceUpdateArgs).Scan(&finalBalance)
		if err != nil {
			return err
		}

		// Вставка записи в "withdraw"
		const stmtWithdrawInsert = `
		INSERT INTO "withdraw" 
		(id, sum, user_id, processed_at)
		VALUES (@id, @sum, @userID, @processedAt)
		RETURNING id`
		withdrawInsertArgs := pgx.NamedArgs{
			"id":          withdraw.ID,
			"sum":         withdraw.Sum,
			"userID":      withdraw.UserID,
			"processedAt": time.Now(),
		}
		var withdrawID string
		err = tx.QueryRow(ctx, stmtWithdrawInsert, withdrawInsertArgs).Scan(&withdrawID)
		if err != nil {
			return err
		}

		return nil
	})
}
func (b *balanceRepository) Withdrawals(ctx context.Context, userID int64) ([]domain.Withdraw, error) {
	var withdrawals []domain.Withdraw
	stmt := `select id,sum, processed_at from "withdraw" where user_id = @userID`
	args := pgx.NamedArgs{
		"userID": userID,
	}
	rows, err := b.database.Pool.Query(ctx, stmt, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.Withdraw{}, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var withdraw domain.Withdraw
		err = rows.Scan(&withdraw.ID, &withdraw.Sum, &withdraw.ProcessedAt)
		if err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, withdraw)
	}
	return withdrawals, nil
}
