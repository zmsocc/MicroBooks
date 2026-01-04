//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	repository2 "github.com/zmsocc/practice/webook/interactive/repository"
	cache2 "github.com/zmsocc/practice/webook/interactive/repository/cache"
	dao2 "github.com/zmsocc/practice/webook/interactive/repository/dao"
	service2 "github.com/zmsocc/practice/webook/interactive/service"
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

var thirdProvider = wire.NewSet(
	InitRedis,
	NewSyncProducer,
	InitKafka,
	InitTestDB, InitLog,
)

var userSvcProvider = wire.NewSet(
	dao.NewUserDAO,
	cache.NewUserCache,
	repository.NewUserRepository,
	service.NewUserService,
)

var interactiveSvcProvider = wire.NewSet(
	service2.NewInteractiveService,
	repository2.NewInteractiveRepository,
	dao2.NewInteractiveDAO,
	cache2.NewRedisInteractiveCache,
)

var articlSvcProvider = wire.NewSet(
	articles.NewArticleDao,
	cache.NewArticleCache,
	articles2.NewArticleRepository,
	service.NewArticleService)

func InitWebServer() *gin.Engine {
	wire.Build(
		thirdProvider,
		userSvcProvider,
		articlSvcProvider,
		interactiveSvcProvider,
		article.NewKafkaProducer,
		cache.NewCodeCache,
		repository.NewCodeRepository,
		// service 部分
		// 集成测试我们显式指定使用内存实现
		ioc.InitSMSService,

		// 指定啥也不干的 wechat service
		//InitPhantomWechatService,
		service.NewCodeService,
		// handler 部分
		web.NewUserHandler,
		//web.NewOAuth2WechatHandler,
		web.NewArticleHandler,
		ijwt.NewRedisJWTHandler,

		// gin 的中间件
		ioc.InitMiddlewares,

		// Web 服务器
		ioc.InitWebServer,
	)
	// 随便返回一个
	return gin.Default()
}

func InitArticleHandler(dao articles.ArticleDAO) *web.ArticleHandler {
	wire.Build(thirdProvider,
		//userSvcProvider,
		interactiveSvcProvider,
		cache.NewArticleCache,
		//ioc.InitIntrGRPCClient,
		//wire.InterfaceValue(new(articles.ArticleDAO), dao),
		article.NewKafkaProducer,
		articles2.NewArticleRepository,
		service.NewArticleService,
		web.NewArticleHandler)
	return new(web.ArticleHandler)
}

func InitUserSvc() service.UserService {
	wire.Build(thirdProvider, userSvcProvider)
	return service.NewUserService(nil)
}

func InitJwtHdl() ijwt.Handler {
	wire.Build(thirdProvider, ijwt.NewRedisJWTHandler)
	return ijwt.NewRedisJWTHandler(nil)
}
