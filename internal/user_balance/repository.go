package userbalance

import (
	"context"
	"database/sql"

	"github.com/citadel-corp/paimon-bank/internal/common/db"
)

type Repository interface {
	RecordBalance(ctx context.Context, payload CreateUserBalancePayload) error
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

// Create implements Repository.
func (d *dbRepository) RecordBalance(ctx context.Context, payload CreateUserBalancePayload) error {
	err := d.db.StartTx(ctx, func(*sql.Tx) error {
		// upsert balance
		upsertBalanceQuery := `
			INSERT INTO user_balance (
				balance, currency, user_id
			) VALUES (
				$1, $2, $3
			)
			ON CONFLICT ON CONSTRAINT user_balance_user_id_currency_unique
			DO UPDATE 
				SET balance = user_balance.balance + $1;
		`
		_, err := d.db.DB().ExecContext(ctx, upsertBalanceQuery, payload.AddedBalance, payload.Currency, payload.UserID)
		if err != nil {
			return err
		}

		// insert into transactions
		createTransactionQuery := `
			INSERT INTO user_transactions (
				user_id, amount, currency, bank_account_number, bank_name, image_url
			) VALUES (
				$1, $2, $3, $4, $5, $6
			)
		`
		_, err = d.db.DB().ExecContext(ctx, createTransactionQuery, payload.UserID, payload.AddedBalance, payload.Currency, payload.SenderBankAccountNumber, payload.SenderBankName, payload.TransferProofImg)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
