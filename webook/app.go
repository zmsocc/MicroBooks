package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zmsocc/practice/webook/internal/event"
)

type App struct {
	web       *gin.Engine
	consumers []event.Consumer
}
