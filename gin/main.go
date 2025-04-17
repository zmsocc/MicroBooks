package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := gin.Default()
	// 静态路由
	//server.GET("/user", func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, "hello world!!!")
	//})
	// 参数路由
	//server.GET("/hello/:name", func(ctx *gin.Context) {
	//	name := ctx.Param("name")
	//	ctx.String(http.StatusOK, "hello,"+name)
	//})
	// 通配符匹配
	//server.GET("/hello/*html", func(ctx *gin.Context) {
	//	path := ctx.Param("html")
	//	ctx.String(http.StatusOK, "匹配上的值是 %s", path)
	//})
	// 查询参数
	server.GET("/hello", func(ctx *gin.Context) {
		id := ctx.Query("id")
		ctx.String(http.StatusOK, "你传过来的 ID 是 %s", id)
	})
	server.Run(":8080")
}
