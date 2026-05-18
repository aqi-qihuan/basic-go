package main

import (
	"basic-go/lmbook/pkg/ginx"
	"basic-go/lmbook/pkg/grpcx"
	"basic-go/lmbook/pkg/saramax"
)

type App struct {
	consumers   []saramax.Consumer
	server      *grpcx.Server
	adminServer *ginx.Server
}
