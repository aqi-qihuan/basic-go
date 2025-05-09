package service

import (
	"basic-go/lmbook/account/domain"
	"basic-go/lmbook/account/repository"
	"context"
)

type accountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

func (a *accountService) Credit(ctx context.Context, cr domain.Credit) error {
	return a.repo.AddCredit(ctx, cr)
}
