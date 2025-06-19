package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zmsocc/practice/webook/internal/web/ijwt"
	"net/http"
)

func main() {
	if err := ijwt.InitPrivateKey(); err != nil {
		panic(err) // 或优雅处理
	}
	//server := gin.Default()
	//server := InitWebServer()
	//server.Run(":8080")

	initPrometheus()

	app := InitWebServer()
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	server := app.web
	//server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好，你来了")
	})
	server.Run(":8080")
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
