package repository

import (
	"basic-go/lmbook/internal/domain"
	"context"
)

type HistoryRecordRepository interface {
	AddRecord(ctx context.Context, record domain.HistoryRecord) error
}
