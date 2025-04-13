package main

import (
	"basic-go/lmbook/internal/events"
	"basic-go/lmbook/pkg/ginx"
	"basic-go/lmbook/pkg/grpcx"
)

type App struct {
	consumers   []events.Consumer
	server      *grpcx.Server
	adminServer *ginx.Server
}
