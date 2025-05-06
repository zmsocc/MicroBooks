package main

import "github.com/zmsocc/practice/webook/internal/web/ijwt"

func main() {
	if err := ijwt.InitPrivateKey(); err != nil {
		panic(err) // 或优雅处理
	}
	//server := gin.Default()
	server := InitWebServer()
	server.Run(":8080")
}
