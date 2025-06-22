package redisx

import (
	"context"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"net"
	"strconv"
	"time"
)

type PrometheusHook struct {
	vector *prometheus.SummaryVec
}

func NewPrometheusHook(opt prometheus.SummaryOpts) *PrometheusHook {
	vector := prometheus.NewSummaryVec(opt, []string{"cmd", "key_exist"})
	err := prometheus.Register(vector)
	if err != nil {
		return nil
	}
	return &PrometheusHook{
		vector: vector,
	}
}

func (p *PrometheusHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		// 相当于，你这里啥也不干
		return next(ctx, network, addr)
	}
}

func (p *PrometheusHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		// 在 Redis 执行之前
		startTime := time.Now()
		var err error
		defer func() {
			duration := time.Since(startTime).Milliseconds()
			//biz := ctx.Value("biz")
			keyExist := errors.Is(redis.Nil, err)
			p.vector.WithLabelValues(
				cmd.Name(),
				//biz.(string),
				strconv.FormatBool(keyExist),
			).Observe(float64(duration))
		}()
		// 这个会最终发送命令到 Redis shang
		err = next(ctx, cmd)
		if err != nil {
			return err
		}
		// 在 Redis 执行之后
		return nil
	}
}

func (p *PrometheusHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	//TODO implement me
	panic("implement me")
}
