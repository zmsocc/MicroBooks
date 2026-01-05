//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/zmsocc/practice/webook/interactive/events"
	"github.com/zmsocc/practice/webook/interactive/grpc"
	"github.com/zmsocc/practice/webook/interactive/ioc"
	"github.com/zmsocc/practice/webook/interactive/repository"
	"github.com/zmsocc/practice/webook/interactive/repository/cache"
	"github.com/zmsocc/practice/webook/interactive/repository/dao"
	"github.com/zmsocc/practice/webook/interactive/service"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitKafka,
	ioc.InitLogger,
)

var interactiveSvcProvider = wire.NewSet(
	service.NewInteractiveService,
	repository.NewInteractiveRepository,
	dao.NewInteractiveDAO,
	cache.NewRedisInteractiveCache,
)

func InitApp() *App {
	wire.Build(
		thirdPartySet,
		interactiveSvcProvider,
		grpc.NewInteractiveServiceServer,
		events.NewInteractiveReadEventConsumer,
		ioc.NewConsumers,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
