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
	SuccessCreateBalance = Response{Code: 200, Message: "Balance added successfully"}
)
