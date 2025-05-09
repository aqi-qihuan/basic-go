package repository

import (
	"basic-go/lmbook/account/domain"
	"context"
)

type AccountRepository interface {
	AddCredit(ctx context.Context, c domain.Credit) error
}
