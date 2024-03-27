package userbalance

import (
	"context"
	"database/sql"
	"errors"
)

type Service interface {
	Create(ctx context.Context, req CreateUserBalancePayload) Response
	CreateTransaction(ctx context.Context, req CreateTransactionPayload) Response
	List(ctx context.Context, req ListUserBalancePayload) Response
}

type userBalanceService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &userBalanceService{repository: repository}
}

func (s *userBalanceService) Create(ctx context.Context, req CreateUserBalancePayload) Response {
	err := s.repository.RecordBalance(ctx, req)
	if err != nil {
		resp := ErrorInternal
		resp.Error = err.Error()
		return resp
	}

	return SuccessCreateBalance
}

// CreateTransaction implements Service.
func (s *userBalanceService) CreateTransaction(ctx context.Context, req CreateTransactionPayload) Response {
	err := s.repository.RecordTransaction(ctx, req)
	if errors.Is(err, ErrNotEnoughBalance) {
		resp := ErrorBadRequest
		resp.Error = err.Error()
		return resp
	}
	if err != nil {
		resp := ErrorInternal
		resp.Error = err.Error()
		return resp
	}

	return SuccessCreateTransaction
}

func (s *userBalanceService) List(ctx context.Context, req ListUserBalancePayload) Response {
	var resp Response

	result, err := s.repository.FindByUserID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return SuccessListBalance
		}
		resp = ErrorInternal
		resp.Error = err.Error()
		return resp
	}

	resp = SuccessListBalance
	resp.Data = result

	return resp
}
