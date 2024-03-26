package userbalance

import "context"

type Service interface {
	Create(ctx context.Context, req CreateUserBalancePayload) Response
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
