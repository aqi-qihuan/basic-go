package client

import (
	intrv1 "basic-go/lmbook/api/proto/gen/intr/v1"
	"basic-go/lmbook/interactive/domain"
	"basic-go/lmbook/interactive/service"
	"context"
	"google.golang.org/grpc"
)

type InteractiveLocalAdapter struct {
	svc service.InteractiveService
}

func NewInteractiveServiceAdapter(svc service.InteractiveService) *InteractiveLocalAdapter {
	return &InteractiveLocalAdapter{svc: svc}
}

func (i *InteractiveLocalAdapter) IncrReadCnt(ctx context.Context, in *intrv1.IncrReadCntRequest, opts ...grpc.CallOption) (*intrv1.IncrReadCntResponse, error) {
	err := i.svc.IncrReadCnt(ctx, in.GetBiz(), in.GetBizId())
	return &intrv1.IncrReadCntResponse{}, err
}

func (i *InteractiveLocalAdapter) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	err := i.svc.Like(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	return &intrv1.LikeResponse{}, err
}

func (i *InteractiveLocalAdapter) CancelLike(ctx context.Context, in *intrv1.CancelLikeRequest, opts ...grpc.CallOption) (*intrv1.CancelLikeResponse, error) {
	err := i.svc.CancelLike(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	return &intrv1.CancelLikeResponse{}, err
}

func (i *InteractiveLocalAdapter) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	err := i.svc.Collect(ctx, in.GetBiz(), in.GetBizId(), in.GetCid(), in.GetUid())
	return &intrv1.CollectResponse{}, err
}

func (i *InteractiveLocalAdapter) Get(ctx context.Context, in *intrv1.GetRequest, opts ...grpc.CallOption) (*intrv1.GetResponse, error) {
	res, err := i.svc.Get(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	if err != nil {
		return nil, err
	}
	return &intrv1.GetResponse{
		Intr: i.toDTO(res),
	}, nil
}

func (i *InteractiveLocalAdapter) GetByIds(ctx context.Context, in *intrv1.GetByIdsRequest, opts ...grpc.CallOption) (*intrv1.GetByIdsResponse, error) {
	if len(in.Ids) == 0 {
		return &intrv1.GetByIdsResponse{}, nil
	}
	data, err := i.svc.GetByIds(ctx, in.GetBiz(), in.GetIds())
	if err != nil {
		return nil, err
	}

	res := make(map[int64]*intrv1.Interactive, len(data))
	for k, v := range data {
		res[k] = i.toDTO(v)
	}
	return &intrv1.GetByIdsResponse{
		Intrs: res,
	}, nil
}

func (i *InteractiveLocalAdapter) CancelCollect(ctx context.Context, in *intrv1.CancelCollectRequest, opts ...grpc.CallOption) (*intrv1.CancelCollectResponse, error) {
	err := i.svc.CancelCollect(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	return &intrv1.CancelCollectResponse{}, err
}

func (i *InteractiveLocalAdapter) GetCollections(ctx context.Context, in *intrv1.GetCollectionsRequest, opts ...grpc.CallOption) (*intrv1.GetCollectionsResponse, error) {
	bizIds, err := i.svc.GetCollections(ctx, in.GetUid(), in.GetBiz(),
		int(in.GetOffset()), int(in.GetLimit()))
	if err != nil {
		return nil, err
	}
	items := make([]*intrv1.CollectionItem, len(bizIds))
	for idx, bizId := range bizIds {
		items[idx] = &intrv1.CollectionItem{
			BizId: bizId,
		}
	}
	return &intrv1.GetCollectionsResponse{
		Items: items,
		Total: int64(len(bizIds)),
	}, nil
}

func (i *InteractiveLocalAdapter) toDTO(intr domain.Interactive) *intrv1.Interactive {
	return &intrv1.Interactive{
		Biz:        intr.Biz,
		BizId:      intr.BizId,
		ReadCnt:    intr.ReadCnt,
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		Liked:      intr.Liked,
		Collected:  intr.Collected,
	}
}
