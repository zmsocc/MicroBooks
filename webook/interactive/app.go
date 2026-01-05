package main

import (
	"github.com/zmsocc/practice/webook/pkg/grpcx"
	"github.com/zmsocc/practice/webook/pkg/saramax"
)

type App struct {
	server    *grpcx.Server
	consumers []saramax.Consumer
}
