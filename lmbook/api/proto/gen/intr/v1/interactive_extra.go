package intrv1

// 本文件包含手动添加的 message 类型
// 对应 interactive.proto 中新增的 CancelCollect/GetCollections 相关 message
// 等 protoc 环境就绪后可重新生成并删除此文件

// CancelCollectRequest 取消收藏请求
type CancelCollectRequest struct {
	Biz   string `protobuf:"bytes,1,opt,name=biz,proto3" json:"biz,omitempty"`
	BizId int64  `protobuf:"varint,2,opt,name=biz_id,json=bizId,proto3" json:"biz_id,omitempty"`
	Uid   int64  `protobuf:"varint,3,opt,name=uid,proto3" json:"uid,omitempty"`
}

func (x *CancelCollectRequest) Reset()         { *x = CancelCollectRequest{} }
func (x *CancelCollectRequest) String() string  { return x.Biz }
func (x *CancelCollectRequest) ProtoMessage()   {}
func (x *CancelCollectRequest) GetBiz() string   { if x != nil { return x.Biz }; return "" }
func (x *CancelCollectRequest) GetBizId() int64  { if x != nil { return x.BizId }; return 0 }
func (x *CancelCollectRequest) GetUid() int64    { if x != nil { return x.Uid }; return 0 }

// CancelCollectResponse 取消收藏响应
type CancelCollectResponse struct{}

func (x *CancelCollectResponse) Reset()        { *x = CancelCollectResponse{} }
func (x *CancelCollectResponse) String() string { return "" }
func (x *CancelCollectResponse) ProtoMessage()  {}

// GetCollectionsRequest 获取收藏列表请求
type GetCollectionsRequest struct {
	Uid    int64  `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Biz    string `protobuf:"bytes,2,opt,name=biz,proto3" json:"biz,omitempty"`
	Offset int32  `protobuf:"varint,3,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit  int32  `protobuf:"varint,4,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *GetCollectionsRequest) Reset()         { *x = GetCollectionsRequest{} }
func (x *GetCollectionsRequest) String() string  { return "" }
func (x *GetCollectionsRequest) ProtoMessage()   {}
func (x *GetCollectionsRequest) GetUid() int64    { if x != nil { return x.Uid }; return 0 }
func (x *GetCollectionsRequest) GetBiz() string   { if x != nil { return x.Biz }; return "" }
func (x *GetCollectionsRequest) GetOffset() int32 { if x != nil { return x.Offset }; return 0 }
func (x *GetCollectionsRequest) GetLimit() int32  { if x != nil { return x.Limit }; return 0 }

// GetCollectionsResponse 获取收藏列表响应
type GetCollectionsResponse struct {
	Items []*CollectionItem `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	Total int64             `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
}

func (x *GetCollectionsResponse) Reset()         { *x = GetCollectionsResponse{} }
func (x *GetCollectionsResponse) String() string  { return "" }
func (x *GetCollectionsResponse) ProtoMessage()   {}
func (x *GetCollectionsResponse) GetItems() []*CollectionItem { if x != nil { return x.Items }; return nil }
func (x *GetCollectionsResponse) GetTotal() int64 { if x != nil { return x.Total }; return 0 }

// CollectionItem 收藏项
type CollectionItem struct {
	BizId int64 `protobuf:"varint,1,opt,name=biz_id,json=bizId,proto3" json:"biz_id,omitempty"`
	Ctime int64 `protobuf:"varint,2,opt,name=ctime,proto3" json:"ctime,omitempty"`
}

func (x *CollectionItem) Reset()         { *x = CollectionItem{} }
func (x *CollectionItem) String() string  { return "" }
func (x *CollectionItem) ProtoMessage()   {}
func (x *CollectionItem) GetBizId() int64 { if x != nil { return x.BizId }; return 0 }
func (x *CollectionItem) GetCtime() int64 { if x != nil { return x.Ctime }; return 0 }
