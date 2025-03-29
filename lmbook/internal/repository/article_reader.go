package repository

import (
	"basic-go/lmbook/internal/domain"
	"context"
)

type ArticleReaderRepository interface {
	Save(ctx context.Context, art domain.Article) error
}
