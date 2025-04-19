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
	tx, err := b.database.Pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	stmtBalance := `
		select balance from balance where userID=@userID 
		`
	balanceArgs := pgx.NamedArgs{
		"userID": userID,
	}
	var balance domain.Balance
	err = tx.QueryRow(ctx, stmtBalance, balanceArgs).Scan(&balance.Balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			stmtBalance = `INSERT INTO "balance" (balance, user_id,created_at,updated_at) VALUES(@balance, @userID, @createdAt, @updatedAt) returning id
			`
			balanceArgs = pgx.NamedArgs{
				"balance":   accrual,
				"userID":    userID,
				"createdAt": time.Now(),
				"updatedAt": time.Now(),
			}
			var id int64
			err = tx.QueryRow(ctx, stmtBalance, balanceArgs).Scan(&id)
			if err != nil {
				return err
			}
		}
		return err
	}
	stmtBalance = `UPDATE "balance" 
	SET balance = balance + @accrual, updated_at = @updatedAt
	WHERE user_id = @userID
	  returning balance
			`
	balanceArgs = pgx.NamedArgs{
		"balance":   accrual,
		"userID":    userID,
		"updatedAt": time.Now(),
	}
	var id int64
	err = tx.QueryRow(ctx, stmtBalance, balanceArgs).Scan(&id)
	if err != nil {
		return nil
	}

	return nil
}

func (b *balanceRepository) Withdraw(ctx context.Context, withdraw domain.Withdraw) error {
	tx, err := b.database.Pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	var withdrawAlreadyExist string
	stmtWithdrawExist := `
	SELECT id from "withdraw" where id=@orderID;
	`
	withdrawExistsArgs := pgx.NamedArgs{
		"orderID": withdraw.ID,
	}
	err = tx.QueryRow(ctx, stmtWithdrawExist, withdrawExistsArgs).Scan(&withdrawAlreadyExist)

	if err == nil {
		return domain.WithdrawAlreadyUsedError
	}

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	var balance float64
	stmtBalance := `select balance.balance from "balance" where user_id=@userID`
	balanceArgs := pgx.NamedArgs{
		"userID": withdraw.UserId,
	}
	err = tx.QueryRow(ctx, stmtBalance, balanceArgs).Scan(&balance)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if errors.Is(err, pgx.ErrNoRows) || balance < withdraw.Sum {
		return domain.InsufficientFundsError
	}
	stmtBalance = `UPDATE "balance" 
	SET balance = balance - @accrual, updated_at = @updatedAt
	WHERE user_id = @userID
	  returning balance`
	balanceArgs = pgx.NamedArgs{
		"userID":    withdraw.UserId,
		"updatedAt": time.Now(),
		"accrual":   withdraw.Sum,
	}
	var finalBalance string
	err = tx.QueryRow(ctx, stmtBalance, balanceArgs).Scan(&finalBalance)
	if err != nil {
		return err
	}

	stmtWithdraw := `
	INSERT INTO "withdraw" 
	(id,sum,user_id,processed_at )
	VALUES (@id,@sum,@userID,@processedAt) 
	returning id;
	`
	args := pgx.NamedArgs{
		"id":          withdraw.ID,
		"sum":         withdraw.Sum,
		"userID":      withdraw.UserId,
		"processedAt": time.Now(),
	}
	var withdrawID string
	err = tx.QueryRow(ctx, stmtWithdraw, args).Scan(&withdrawID)
	return err
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
