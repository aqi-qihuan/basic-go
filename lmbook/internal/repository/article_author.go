package repository

import (
	"basic-go/lmbook/internal/domain"
	"context"
)

type ArticleAuthorRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
}
