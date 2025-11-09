package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type JobDAO interface {
	Preempt(ctx context.Context) (Job, error)
	Release(ctx context.Context, id int64) error
	UpdateUtime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, next time.Time) error
	Stop(ctx context.Context, id int64) error
}

type GORMJobDAO struct {
	db *gorm.DB
}

func (g *GORMJobDAO) Preempt(ctx context.Context) (Job, error) {
	db := g.db.WithContext(ctx)
	for {
		now := time.Now()
		var j Job
		err := db.WithContext(ctx).Model(&Job{}).
			Where("status = ? AND next_time <= ?", jobStatusWaiting, now).First(&j).Error
		// 找到了，可以被抢占的
		// 找到之后要干嘛？要抢占
		if err != nil {
			return Job{}, err
		}
		// 两个 goroutine 都拿到 id = 1 的数据
		// 能不能用 utime？
		// 乐观锁，CAS 操作，compare AND Swap
		// 可以用乐观锁取代 FOR UPDATE
		res := db.WithContext(ctx).Model(&Job{}).
			Where("id = ? AND version = ?", j.Id, j.Version).
			Updates(map[string]any{
				"status":  jobStatusRunning,
				"utime":   now,
				"version": j.Version + 1,
			})
		if res.Error != nil {
			return Job{}, err
		}
		if res.RowsAffected == 0 {
			continue
		}
		return j, nil
	}
}

func (g *GORMJobDAO) Release(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", id).Updates(map[string]any{
		"status": jobStatusWaiting,
		"utime":  time.Now().UnixMilli(),
	}).Error
}

func (g *GORMJobDAO) UpdateUtime(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", id).Updates(map[string]any{
		"utime": time.Now().UnixMilli(),
	}).Error
}

func (g *GORMJobDAO) UpdateNextTime(ctx context.Context, id int64, next time.Time) error {
	return g.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", id).Updates(map[string]any{
		"next_time": next.UnixMilli(),
	}).Error
}

func (g *GORMJobDAO) Stop(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", id).Updates(map[string]any{
		"status": jobStatusPaused,
		"utime":  time.Now().UnixMilli(),
	}).Error
}

type Job struct {
	Id       int64 `gorm:"primary_key, autoIncrement"`
	Cfg      string
	Executor string
	Name     string `gorm:"unique"`
	// 哪些任务可以抢，哪些任务已经被人占着，哪些任务永远不会被运行
	// 用状态来标记
	Status int
	// 定时任务，我怎么知道已经到时间了呢？
	// NextTime 下一次被调度的时间
	// next_time <= now 这样一个查询条件
	// and status = 0
	// 要建立索引
	// 更好的应该是 next_tim 和 status 的联合索引
	NextTime int64 `gorm:"index"`
	// cron 表达式
	Cron    string
	Version int
	// 创建时间，毫秒数
	Ctime int64
	// 更新时间，毫秒数
	Utime int64
}

const (
	jobStatusWaiting = iota
	// 已经被抢占
	jobStatusRunning
	// 还可以有别的取值
	// 暂停调度
	jobStatusPaused
)
