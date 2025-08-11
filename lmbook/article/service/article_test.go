package service

import (
	"basic-go/lmbook/article/domain"
	"basic-go/lmbook/article/repository"
	repomocks "basic-go/lmbook/article/repository/mocks"
	"basic-go/lmbook/pkg/logger"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestArticleService_PublishV1(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (
			authorRepo repository.ArticleAuthorRepository,
			readerRepo repository.ArticleReaderRepository)

		art domain.Article

		wantErr error
		wantId  int64
	}{
		{
			name: "新建发表成功",
			mock: func(ctrl *gomock.Controller) (
				authorRepo repository.ArticleAuthorRepository,
				readerRepo repository.ArticleReaderRepository) {
				ar := repomocks.NewMockArticleAuthorRepository(ctrl)
				ar.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)

				rr := repomocks.NewMockArticleReaderRepository(ctrl)
				rr.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      1,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				return ar, rr
			},
			art: domain.Article{
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId: 1,
		},
		{
			name: "修改保存到制作库失败",
			mock: func(ctrl *gomock.Controller) (
				authorRepo repository.ArticleAuthorRepository,
				readerRepo repository.ArticleReaderRepository) {
				ar := repomocks.NewMockArticleAuthorRepository(ctrl)
				ar.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      7,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(errors.New("保存失败"))
				rr := repomocks.NewMockArticleReaderRepository(ctrl)
				return ar, rr
			},
			art: domain.Article{
				Id:      7,
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantErr: errors.New("保存失败"),
		},
		{
			name: "修改保存到线上库失败-重试都失败了",
			mock: func(ctrl *gomock.Controller) (
				authorRepo repository.ArticleAuthorRepository,
				readerRepo repository.ArticleReaderRepository) {
				ar := repomocks.NewMockArticleAuthorRepository(ctrl)
				ar.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      7,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)

				rr := repomocks.NewMockArticleReaderRepository(ctrl)
				// 模拟多次重试均失败
				rr.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      7,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).AnyTimes().Return(errors.New("保存到线上库失败"))
				return ar, rr
			},
			art: domain.Article{
				Id:      7,
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantErr: errors.New("保存到线上库失败"),
		},
		{
			name: "修改保存到线上库失败-重试成功",
			mock: func(ctrl *gomock.Controller) (
				authorRepo repository.ArticleAuthorRepository,
				readerRepo repository.ArticleReaderRepository) {
				ar := repomocks.NewMockArticleAuthorRepository(ctrl)
				ar.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      7,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)

				rr := repomocks.NewMockArticleReaderRepository(ctrl)
				// 第一次失败
				rr.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      7,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(errors.New("保存到线上库失败"))
				// 第二次成功
				rr.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      7,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				return ar, rr
			},
			art: domain.Article{
				Id:      7,
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId: 7,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			authorRepo, readerRepo := tc.mock(ctrl)
			svc := NewArticleServiceV1(authorRepo, readerRepo,
				nil, logger.NewNoOpLogger())
			id, err := svc.PublishV1(context.Background(), tc.art)
			if tc.wantErr != nil {
				// 使用ErrorContains来检查错误消息，因为测试用例中的errors.New创建的是临时错误对象
				assert.ErrorContains(t, err, tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantId, id)
		})
	}
}
