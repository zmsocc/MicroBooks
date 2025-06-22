package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zmsocc/practice/webook/pkg/logger"
	"net/http"
	"strconv"
)

// L 使用包变量
var L logger.Logger

var vector *prometheus.CounterVec

func InitCounter(opt prometheus.CounterOpts) {
	vector = prometheus.NewCounterVec(opt, []string{"code"})
	err := prometheus.Register(vector)
	if err != nil {
		return
	}
	// 你可以考虑使用 code， method， 命中路由， HTTP 状态码
}

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

func WrapToken[C jwt.Claims](fn func(ctx *gin.Context, uc C) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 执行一些东西
		val, ok := ctx.Get("users")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c, ok := val.(C)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 下半段的业务逻辑从哪里来？
		// 我们的业务逻辑有可能要操作 ctx
		// 你要读取 HTTP HEADER
		res, err := fn(ctx, c)
		if err != nil {
			// 开始处理 error， 其实就是记录一下日志
			L.Error("处理业务逻辑出错",
				logger.String("path", ctx.Request.URL.Path),
				logger.String("route", ctx.FullPath()),
				logger.Error(err))
		}
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		ctx.JSON(http.StatusOK, res)
		// 再执行一些东西
	}
}
