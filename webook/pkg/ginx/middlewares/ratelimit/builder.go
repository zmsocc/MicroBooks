package ratelimit

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

type Builder struct {
	prefix   string
	cmd      redis.Cmdable
	interval time.Duration
	// 阈值
	rate int
}

//go:embed slide_widow.lua
var luaScript string

func NewBuilder(cmd redis.Cmdable, interval time.Duration, rate int) *Builder {
	return &Builder{
		prefix:   "ip-limiter",
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (b *Builder) Prefix(prefix string) *Builder {
	b.prefix = prefix
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limited, err := b.limit(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if limited {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}

func (b *Builder) limit(ctx *gin.Context) (bool, error) {
	key := fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP())
	return b.cmd.Eval(ctx, luaScript, []string{key}, b.interval.Milliseconds(), b.rate, time.Now().UnixMilli()).Bool()
}
