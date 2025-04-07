// Code generated by MockGen. DO NOT EDIT.
// Source: ./lmbook/internal/repository/article.go
//
// Generated by this command:
//
//	mockgen -source=./lmbook/internal/repository/article.go -package=repomocks -destination=./lmbook/internal/repository/mocks/article.mock.go
//
// Package repomocks is a generated GoMock package.
package repomocks

import (
	domain "basic-go/lmbook/internal/domain"
	context "context"
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockArticleRepository is a mock of ArticleRepository interface.
type MockArticleRepository struct {
	ctrl     *gomock.Controller
	recorder *MockArticleRepositoryMockRecorder
}

// MockArticleRepositoryMockRecorder is the mock recorder for MockArticleRepository.
type MockArticleRepositoryMockRecorder struct {
	mock *MockArticleRepository
}

// NewMockArticleRepository creates a new mock instance.
func NewMockArticleRepository(ctrl *gomock.Controller) *MockArticleRepository {
	mock := &MockArticleRepository{ctrl: ctrl}
	mock.recorder = &MockArticleRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockArticleRepository) EXPECT() *MockArticleRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, art)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockArticleRepositoryMockRecorder) Create(ctx, art any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockArticleRepository)(nil).Create), ctx, art)
}

// GetByAuthor mocks base method.
func (m *MockArticleRepository) GetByAuthor(ctx context.Context, uid int64, offset, limit int) ([]domain.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByAuthor", ctx, uid, offset, limit)
	ret0, _ := ret[0].([]domain.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByAuthor indicates an expected call of GetByAuthor.
func (mr *MockArticleRepositoryMockRecorder) GetByAuthor(ctx, uid, offset, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByAuthor", reflect.TypeOf((*MockArticleRepository)(nil).GetByAuthor), ctx, uid, offset, limit)
}

// GetById mocks base method.
func (m *MockArticleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", ctx, id)
	ret0, _ := ret[0].(domain.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockArticleRepositoryMockRecorder) GetById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockArticleRepository)(nil).GetById), ctx, id)
}

// GetPubById mocks base method.
func (m *MockArticleRepository) GetPubById(ctx context.Context, id int64) (domain.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPubById", ctx, id)
	ret0, _ := ret[0].(domain.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPubById indicates an expected call of GetPubById.
func (mr *MockArticleRepositoryMockRecorder) GetPubById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPubById", reflect.TypeOf((*MockArticleRepository)(nil).GetPubById), ctx, id)
}

// ListPub mocks base method.
func (m *MockArticleRepository) ListPub(ctx context.Context, start time.Time, offset, limit int) ([]domain.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPub", ctx, start, offset, limit)
	ret0, _ := ret[0].([]domain.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPub indicates an expected call of ListPub.
func (mr *MockArticleRepositoryMockRecorder) ListPub(ctx, start, offset, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPub", reflect.TypeOf((*MockArticleRepository)(nil).ListPub), ctx, start, offset, limit)
}

// Sync mocks base method.
func (m *MockArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sync", ctx, art)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sync indicates an expected call of Sync.
func (mr *MockArticleRepositoryMockRecorder) Sync(ctx, art any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*MockArticleRepository)(nil).Sync), ctx, art)
}

// SyncStatus mocks base method.
func (m *MockArticleRepository) SyncStatus(ctx context.Context, uid, id int64, status domain.ArticleStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncStatus", ctx, uid, id, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// SyncStatus indicates an expected call of SyncStatus.
func (mr *MockArticleRepositoryMockRecorder) SyncStatus(ctx, uid, id, status any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncStatus", reflect.TypeOf((*MockArticleRepository)(nil).SyncStatus), ctx, uid, id, status)
}

// Update mocks base method.
func (m *MockArticleRepository) Update(ctx context.Context, art domain.Article) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, art)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockArticleRepositoryMockRecorder) Update(ctx, art any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockArticleRepository)(nil).Update), ctx, art)
}
