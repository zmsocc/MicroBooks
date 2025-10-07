package service

import (
	"context"
	"github.com/zmsocc/practice/webook/internal/domain"
	"github.com/zmsocc/practice/webook/internal/repository"
	"github.com/zmsocc/practice/webook/pkg/logger"
	"time"
)

type JobService interface {
	// Preempt 抢占
	Preempt(ctx context.Context) (domain.Job, error)
	Release(ctx context.Context, id int64) error
	ResetNextTime(ctx context.Context, j domain.Job) error
}

type cronJobService struct {
	repo            repository.JobRepository
	refreshInterval time.Duration
	l               logger.Logger
}

func (p *cronJobService) Preempt(ctx context.Context) (domain.Job, error) {
	j, err := p.repo.Preempt(ctx)
	// 续约？
	//ch := make(chan struct{})
	//go func() {
	//	ticker := time.NewTicker(p.refreshInterval)
	//	for {
	//		select {
	//		case <-ticker.C:
	//			// 在这里续约
	//			p.refresh(j.Id)
	//		case <-ch:
	//			// 结束
	//			return
	//		}
	//	}
	//}()
	// 抢占之后，一直抢占着吗？
	// 你要考虑一个释放的问题

	ticker := time.NewTicker(p.refreshInterval)
	go func() {
		for range ticker.C {
			p.refresh(j.Id)
		}
	}()

	j.CancelFunc = func() error {
		//close(ch)
		// 自己在这里释放掉
		ticker.Stop()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		return p.repo.Release(ctx, j.Id)
	}
	return j, err
}

func (p *cronJobService) Release(ctx context.Context, id int64) error {
	return p.repo.Release(ctx, id)
}

func (p *cronJobService) ResetNextTime(ctx context.Context, j domain.Job) error {
	next := j.NextTime()
	if next.IsZero() {
		// 没有下一次
		return p.repo.Stop(ctx, j.Id)
	}
	return p.repo.UpdateNextTime(ctx, j.Id, next)
}

func (p *cronJobService) refresh(id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 续约怎么个续法
	// 更新一下更新时间就可以
	// 比如说我们的续约失败逻辑就是：处于 running 状态，但是更新时间在三分钟以前
	err := p.repo.UpdateUtime(ctx, id)
	if err != nil {
		// 可以考虑立刻重试
		p.l.Error("续约失败", logger.Error(err), logger.Int64("jid", id))
	}
}
