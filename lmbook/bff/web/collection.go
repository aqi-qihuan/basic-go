package web

import (
	intrv1 "basic-go/lmbook/api/proto/gen/intr/v1"
	"basic-go/lmbook/bff/web/jwt"
	"basic-go/lmbook/pkg/ginx"
	"github.com/gin-gonic/gin"
)

var _ handler = (*CollectionHandler)(nil)

type CollectionHandler struct {
	intrSvc intrv1.InteractiveServiceClient
	biz     string
}

func NewCollectionHandler(intrSvc intrv1.InteractiveServiceClient) *CollectionHandler {
	return &CollectionHandler{
		intrSvc: intrSvc,
		biz:     "article",
	}
}

func (h *CollectionHandler) RegisterRoutes(s *gin.Engine) {
	g := s.Group("/collections")
	g.POST("/list", ginx.WrapClaimsAndReq(h.List))
	g.POST("/cancel", ginx.WrapClaimsAndReq(h.Cancel))
}

type CollectionListReq struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type CollectionCancelReq struct {
	BizId int64 `json:"bizId"`
}

// List 获取用户收藏列表
func (h *CollectionHandler) List(ctx *gin.Context,
	req CollectionListReq, uc jwt.UserClaims) (ginx.Result, error) {
	if req.Limit <= 0 || req.Limit > 50 {
		req.Limit = 20
	}

	resp, err := h.intrSvc.GetCollections(ctx, &intrv1.GetCollectionsRequest{
		Uid:    uc.Id,
		Biz:    h.biz,
		Offset: int32(req.Offset),
		Limit:  int32(req.Limit),
	})
	if err != nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}

	return ginx.Result{
		Data: map[string]any{
			"items": resp.Items,
			"total": resp.Total,
		},
	}, nil
}

// Cancel 取消收藏
func (h *CollectionHandler) Cancel(ctx *gin.Context,
	req CollectionCancelReq, uc jwt.UserClaims) (ginx.Result, error) {
	_, err := h.intrSvc.CancelCollect(ctx, &intrv1.CancelCollectRequest{
		Biz:   h.biz,
		BizId: req.BizId,
		Uid:   uc.Id,
	})
	if err != nil {
		return ginx.Result{Code: 5, Msg: "取消收藏失败"}, err
	}
	return ginx.Result{Msg: "OK"}, nil
}
