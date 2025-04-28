package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/zmsocc/practice/webook/internal/service/sms"
	"github.com/zmsocc/practice/webook/internal/service/sms/memory"
)

func InitSMSService(cmd redis.Cmdable) sms.Service {
	// tencent.InitSmsTencentService()
	return memory.NewService()
}
