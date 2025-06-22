package ioc

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/zmsocc/practice/webook/internal/repository/cache"
	"github.com/zmsocc/practice/webook/pkg/redisx"
)

// InitUserHook 配合 PrometheusHook 使用
func InitUserHook(client *redis.Client) cache.UserCache {
	client.AddHook(redisx.NewPrometheusHook(
		prometheus.SummaryOpts{
			Namespace: "geekbang_daming",
			Subsystem: "webook",
			Name:      "sms_resp_time",
			Help:      "统计 SMS 服务的性能数据",
			ConstLabels: map[string]string{
				"biz": "user",
			},
		}))
	panic("你别调用")
}
