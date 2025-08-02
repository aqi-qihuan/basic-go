package service

import (
	"basic-go/lmbook/feed/domain"
	"context"
)

type FeedService interface {
	CreateFeedEvent(ctx context.Context, feed domain.FeedEvent) error
	GetFeedEventList(ctx context.Context, uid, timestamp, limit int64) ([]domain.FeedEvent, error)
}

// Handler 具体业务处理逻辑
type Handler interface {
	CreateFeedEvent(ctx context.Context, feed domain.FeedEvent) error
	FindFeedEvents(ctx context.Context, uid, timestamp, limit int64) ([]domain.FeedEvent, error)
}
