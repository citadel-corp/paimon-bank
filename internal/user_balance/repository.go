package userbalance

import (
	"context"
	"database/sql"
	"errors"

	"github.com/citadel-corp/paimon-bank/internal/common/db"
	"github.com/citadel-corp/paimon-bank/internal/common/id"
	"github.com/citadel-corp/paimon-bank/internal/common/response"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	RecordBalance(ctx context.Context, payload CreateUserBalancePayload) error
	RecordTransaction(ctx context.Context, payload CreateTransactionPayload) error
	FindByUserID(ctx context.Context, userID string) ([]UserBalanceResponse, error)
	ListTransactions(ctx context.Context, payload ListUserTransactionPayload) ([]UserTransaction, *response.Pagination, error)
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
			RETURNING id
		`
		row := tx.QueryRowContext(ctx, updateBalanceQuery, payload.Balances, payload.UserID, payload.FromCurrency)
		var userBalanceID uint64
		err := row.Scan(&userBalanceID)
		var pgErr *pgconn.PgError
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrNoCurrencyOrUserRecorded
			}
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

		if userBalanceID == 0 {
			return ErrNoCurrencyOrUserRecorded
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
	response := []UserBalanceResponse{}

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

// ListTransactions implements Repository.
func (d *dbRepository) ListTransactions(ctx context.Context, payload ListUserTransactionPayload) ([]UserTransaction, *response.Pagination, error) {
	var resp []UserTransaction
	pagination := &response.Pagination{
		Limit:  payload.Limit,
		Offset: payload.Offset,
	}

	selectQuery := `
		SELECT COUNT(*) OVER() AS total_count, *
		FROM user_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
		OFFSET $3
	`

	rows, err := d.db.DB().QueryContext(ctx, selectQuery, payload.UserID, payload.Limit, payload.Offset)
	if err != nil {
		return nil, nil, err
	}

	for rows.Next() {
		var ut UserTransaction
		err = rows.Scan(&pagination.Total, &ut.TransactionID, &ut.UserID, &ut.Amount, &ut.Currency, &ut.BankAccountNumber, &ut.BankName, &ut.ImageURL, &ut.CreatedAt)
		if err != nil {
			return nil, nil, err
		}

		resp = append(resp, ut)
	}

	return resp, pagination, nil
}
