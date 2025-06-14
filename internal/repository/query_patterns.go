package repository

// order stmnts
const (
	GetUnhandledOrdersGetOrder = `SELECT id, status, accrual, user_id, uploaded_at FROM "order" WHERE status IN ('%s', '%s')`
	SaveGetOrderByID           = `SELECT id, status, accrual, user_id, uploaded_at 
	FROM "order" 
	WHERE id = @orderId`
	SaveInsertNewOrder = `INSERT INTO "order" 
	(id, status, accrual, user_id, uploaded_at) 
	VALUES (@id, @status, @accrual, @userID, @uploadedAt) 
	RETURNING id`
	GetAllSelectAllOrderByUserID = `SELECT id, status, accrual, user_id, uploaded_at FROM "order" WHERE user_id = @userID`
	UpdateOrderUpdateOrderByID   = `UPDATE "order" SET status = @status, accrual = @accrual WHERE id = @id`
)

// balance stmnts
const (
	GetBalanceGetBalanceEntityByUserID = `SELECT id,balance,user_id,created_at,updated_at FROM "balance" where user_id=@userID`
	UpdateBalanceGetBalanceByUserID    = `
		SELECT balance FROM "balance" WHERE user_id = @userID`
	UpdateBalanceInsertNewBalanceEntity = `
				INSERT INTO "balance" (balance, user_id, created_at, updated_at)
				VALUES (@balance, @userID, @createdAt, @updatedAt)
				RETURNING id`
	UpdateBalanceUpdateBalance = `
		UPDATE "balance"
		SET balance = balance + @accrual, updated_at = @updatedAt
		WHERE user_id = @userID
		RETURNING balance`
	WithdrawGetWithdrawByUserID = `
		SELECT id FROM "withdraw" WHERE order_id = @orderID`
	WithdrawUpdateBalance = `
		UPDATE "balance" 
		SET balance = balance - @accrual, updated_at = @updatedAt
		WHERE user_id = @userID
		RETURNING balance`
)
