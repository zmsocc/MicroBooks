package main

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zmsocc/practice/webook/internal/web/ijwt"
	"github.com/zmsocc/practice/webook/ioc"
	"net/http"
	"time"
)

func main() {
	if err := ijwt.InitPrivateKey(); err != nil {
		panic(err) // 或优雅处理
	}
	//server := gin.Default()
	//server := InitWebServer()
	//server.Run(":8080")
	initPrometheus()
	closeFunc := ioc.InitOTEL()
	app := InitWebServer()
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	app.cron.Start()
	server := app.web
	//server := gin.Default()
	//server.GET("/hello", func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, "你好，你来了")
	//})
	server.Run(":8080")
	// 一分钟内你要关完退出
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	ctx = app.cron.Stop()
	// 这边可以考虑超时强制退出，防止有些任务，执行特别长的时间
	tm := time.NewTimer(time.Minute * 10)
	select {
	case <-ctx.Done():
	case <-tm.C:
	}
	closeFunc(ctx)
}

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			return
		}
	}()
}
