package userbalance

import "github.com/citadel-corp/paimon-bank/internal/common/response"

type Response struct {
	Code    int
	Message string
	Data    any
	Meta    *response.Pagination
	Error   string
}

var (
	SuccessCreateBalance     = Response{Code: 200, Message: "Balance added successfully"}
	SuccessCreateTransaction = Response{Code: 200, Message: "Transaction successful"}
	Success                  = Response{Code: 200, Message: "success"}
)

type UserBalanceResponse struct {
	Balance  int    `json:"balance"`
	Currency string `json:"currency"`
}

type UserTransactionResponse struct {
	TransactionID    string `json:"transactionId"`
	Balance          int    `json:"balance"`
	Currency         string `json:"currency"`
	TransferProofImg string `json:"transferProofImg"`
	CreatedAt        int64  `json:"createdAt"`
	Source           struct {
		BankAccountNumber string `json:"bankAccountNumber"`
		BankName          string `json:"bankName"`
	} `json:"source"`
}
