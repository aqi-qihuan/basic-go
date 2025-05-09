package service

import (
	"basic-go/lmbook/account/domain"
	"context"
)

type AccountService interface {
	Credit(ctx context.Context, cr domain.Credit) error
}
