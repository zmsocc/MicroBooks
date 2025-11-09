package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/zmsocc/practice/webook/internal/event"
)

type App struct {
	web       *gin.Engine
	consumers []event.Consumer
	cron      *cron.Cron
}
