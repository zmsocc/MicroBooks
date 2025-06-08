package main

import (
	"github.com/gin-gonic/gin"
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

	app := InitWebServer()
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	server := app.web
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好，你来了")
	})
	server.Run(":8080")
}
