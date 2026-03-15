package account

import (
	"context"

	"github.com/Aym-Aymen777/gRPC-GraphQL-microservices/account/types"
	"github.com/Aym-Aymen777/gRPC-GraphQL-microservices/utils"
)

type Service interface {
	CreateAccount(ctx context.Context, username string, email string) (*types.Account, error)
	GetAccountDetails(ctx context.Context, id string) (*types.Account, error)
	GetAccounts(ctx context.Context, skip, limit uint64) ([]*types.Account, error)
}

type accountService struct {
	repo Repository
}

func NewAccountService(repo Repository) Service {
	return &accountService{repo: repo}
}

func (s *accountService) CreateAccount(ctx context.Context, username string, email string) (*types.Account, error) {
	account := &types.Account{
		ID:       utils.GenerateID(),
		Username: username,
		Email:    email,
	}
	err := s.repo.PutAccount(ctx, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *accountService) GetAccountDetails(ctx context.Context, id string) (*types.Account, error) {
	return s.repo.GetAccount(ctx, id)
}

func (s *accountService) GetAccounts(ctx context.Context, skip, limit uint64) ([]*types.Account, error) {
	if limit == 0 {
		limit = 10
	}
	if skip > 0 {
		skip = (skip - 1) * limit
	}
	return s.repo.ListAccounts(ctx, skip, limit)
}
