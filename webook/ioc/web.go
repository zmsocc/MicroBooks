package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zmsocc/practice/webook/internal/web"
	"github.com/zmsocc/practice/webook/internal/web/ijwt"
	"github.com/zmsocc/practice/webook/internal/web/middleware"
	"github.com/zmsocc/practice/webook/pkg/ginx/middlewares/metric"
	"github.com/zmsocc/practice/webook/pkg/ginx/middlewares/ratelimit"
	"strings"
	"time"
)

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"X-Jwt-Token", "x-refresh-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour,
	})
}

func InitMiddlewares(jwtHdl ijwt.Handler, cmd redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		(&metric.MiddlewareBuilder{
			Namespace: "geekbang_daming",
			Subsystem: "webook",
			Name:      "gin_http",
			// 上面三个不能使用连字符-
			Help:       "统计 GIN 的 GTTP 接口",
			InstanceID: "my-instance-1",
		}).Build(),
		middleware.NewLoginJWTMiddlewareBuilder(jwtHdl).
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			IgnorePaths("/users/refresh_token").
			IgnorePaths("/test/metrics").
			Build(),
		ratelimit.NewBuilder(cmd, time.Minute, 100).Build(),
	}
}

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler,
	articleHdl *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	articleHdl.RegisterRoutes(server)
	(&web.ObservabilityHandler{}).RegisterRoutes(server)
	return server
}
