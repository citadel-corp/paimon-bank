package userbalance

import (
	"context"
	"database/sql"
	"errors"

	"github.com/citadel-corp/paimon-bank/internal/common/db"
	"github.com/citadel-corp/paimon-bank/internal/common/id"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	RecordBalance(ctx context.Context, payload CreateUserBalancePayload) error
	RecordTransaction(ctx context.Context, payload CreateTransactionPayload) error
	FindByUserID(ctx context.Context, userID string) ([]UserBalanceResponse, error)
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

// Create implements Repository.
func (d *dbRepository) RecordBalance(ctx context.Context, payload CreateUserBalancePayload) error {
	return d.db.StartTx(ctx, func(tx *sql.Tx) error {
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
		_, err := tx.ExecContext(ctx, upsertBalanceQuery, payload.AddedBalance, payload.Currency, payload.UserID)
		if err != nil {
			return err
		}

		// insert into transactions
		createTransactionQuery := `
			INSERT INTO user_transactions (
				id, user_id, amount, currency, bank_account_number, bank_name, image_url
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7
			)
		`
		_, err = tx.ExecContext(ctx, createTransactionQuery, id.GenerateStringID(16), payload.UserID, payload.AddedBalance, payload.Currency, payload.SenderBankAccountNumber, payload.SenderBankName, payload.TransferProofImg)
		if err != nil {
			return err
		}

		return nil
	})
}

func (d *dbRepository) RecordTransaction(ctx context.Context, payload CreateTransactionPayload) error {
	return d.db.StartTx(ctx, func(tx *sql.Tx) error {
		updateBalanceQuery := `
			UPDATE user_balance
			SET balance = balance - $1
			WHERE user_id = $2 and currency = $3
		`
		_, err := tx.ExecContext(ctx, updateBalanceQuery, payload.Balances, payload.UserID, payload.FromCurrency)
		var pgErr *pgconn.PgError
		if err != nil {
			if errors.As(err, &pgErr) {
				switch pgErr.Code {
				case "23514":
					if pgErr.ConstraintName == "balance_non_negative" {
						return ErrNotEnoughBalance
					}
					return err
				default:
					return err
				}
			}
			return err
		}

		// insert into transactions
		createTransactionQuery := `
			INSERT INTO user_transactions (
				id, user_id, amount, currency, bank_account_number, bank_name
			) VALUES (
				$1, $2, $3, $4, $5, $6
			)
		`
		_, err = tx.ExecContext(ctx, createTransactionQuery, id.GenerateStringID(16), payload.UserID, -payload.Balances, payload.FromCurrency, payload.RecipientBankAccountNumber, payload.RecipientBankName)
		if err != nil {
			return err
		}

		return nil
	})
}

func (d *dbRepository) FindByUserID(ctx context.Context, userID string) ([]UserBalanceResponse, error) {
	var response []UserBalanceResponse

	selectQuery := `
		SELECT balance, currency
		FROM user_balance
		WHERE user_id = $1
		ORDER BY balance desc
	`

	rows, err := d.db.DB().QueryContext(ctx, selectQuery, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ub UserBalanceResponse
		err = rows.Scan(&ub.Balance, &ub.Currency)
		if err != nil {
			return nil, err
		}

		response = append(response, ub)
	}

	return response, nil
}
