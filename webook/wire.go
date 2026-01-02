//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/zmsocc/practice/webook/internal/event/article"
	"github.com/zmsocc/practice/webook/internal/repository"
	articles2 "github.com/zmsocc/practice/webook/internal/repository/articles"
	"github.com/zmsocc/practice/webook/internal/repository/cache"
	"github.com/zmsocc/practice/webook/internal/repository/dao"
	"github.com/zmsocc/practice/webook/internal/repository/dao/articles"
	"github.com/zmsocc/practice/webook/internal/service"
	"github.com/zmsocc/practice/webook/internal/web"
	"github.com/zmsocc/practice/webook/internal/web/ijwt"
	"github.com/zmsocc/practice/webook/ioc"
)

var rankingServiceSet = wire.NewSet(
	repository.NewCachedRankingRepository,
	cache.NewRankingRedisCache,
	service.NewBatchRankingService,
)

var interactiveSvcProvider = wire.NewSet(
	service.NewInteractiveService,
	repository.NewInteractiveRepository,
	dao.NewInteractiveDAO,
	cache.NewRedisInteractiveCache,
)

func InitWebServer() *App {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitLogger,
		ioc.InitKafka,
		ioc.NewConsumers,
		ioc.NewSyncProducer,

		rankingServiceSet,
		interactiveSvcProvider,
		ioc.InitRLockClient,
		ioc.InitJobs,
		ioc.InitRankingJob,

		// consumer
		article.NewInteractiveReadEventBatchConsumer,
		article.NewKafkaProducer,

		// 初始化 DAO
		dao.NewUserDAO,
		articles.NewArticleDao,

		cache.NewUserCache,
		cache.NewCodeCache,
		cache.NewArticleCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,
		articles2.NewArticleRepository,

		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,

		// 直接基于内存实现
		ioc.InitSMSService,

		web.NewUserHandler,
		web.NewArticleHandler,
		ijwt.NewRedisJWTHandler,

		ioc.InitWebServer,
		ioc.InitMiddlewares,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
