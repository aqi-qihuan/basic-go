package web

import "basic-go/lmbook/pkg/ginx"

type FindFeedEventReq struct {
	UID       int64 `json:"uid"`
	Limit     int64 `json:"limit"`
	Timestamp int64 `json:"timestamp"`
}

type CreateFeedEventReq struct {
	Typ string `json:"typ"`
	Ext string `json:"ext"`
}

type Result = ginx.Result
