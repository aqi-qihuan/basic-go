package main

import (
	"basic-go/lmbook/internal/events"
	"basic-go/lmbook/pkg/grpcx"
)

type App struct {
	consumers []events.Consumer
	server    *grpcx.Server
}
