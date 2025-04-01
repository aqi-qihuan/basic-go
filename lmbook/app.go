package main

import (
	"basic-go/lmbook/internal/events"
	"github.com/gin-gonic/gin"
)

type App struct {
	server    *gin.Engine
	consumers []events.Consumer
}
