package userbalance

import "time"

type UserBalance struct {
	ID        uint64     `json:"-"`
	Balance   int        `json:"balance"`
	Currency  string     `json:"currency"`
	UserID    uint64     `json:"user_id"`
	CreatedAt *time.Time `json:"created_at"`
}

type UserTransaction struct {
	TransactionID     string
	UserID            string
	Amount            int
	Currency          string
	BankAccountNumber string
	BankName          string
	ImageURL          *string
	CreatedAt         time.Time
}
