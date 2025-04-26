//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/zmsocc/practice/webook/internal/repository"
	"github.com/zmsocc/practice/webook/internal/repository/cache"
	"github.com/zmsocc/practice/webook/internal/repository/dao"
	"github.com/zmsocc/practice/webook/internal/service"
	"github.com/zmsocc/practice/webook/internal/web"
	"github.com/zmsocc/practice/webook/internal/web/ijwt"
	"github.com/zmsocc/practice/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB,
		ioc.InitRedis,
		// 初始化 DAO
		dao.NewUserDAO,
		cache.NewUserCache,
		cache.NewCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,

		service.NewUserService,
		service.NewCodeService,

		// 直接基于内存实现
		ioc.InitSMSService,

		web.NewUserHandler,
		ijwt.NewRedisJWTHandler,

		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
