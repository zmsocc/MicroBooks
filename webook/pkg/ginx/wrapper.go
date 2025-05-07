package ginx

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Result struct {
	// 业务错误码
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func WrapBody(fn func(ctx *gin.Context) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 要读取 HTTP HEADER
		res, err := fn(ctx)
		if err != nil {
			return
		}
		ctx.JSON(http.StatusOK, res)
	}
}
