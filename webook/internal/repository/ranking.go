package repository

import (
	"context"
	"github.com/zmsocc/practice/webook/internal/domain"
	"github.com/zmsocc/practice/webook/internal/repository/cache"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}

type CachedRankingRepository struct {
	// 使用具体实现，可读性更好，对测试不友好，因为没有面向接口编程
	redis *cache.RankingRedisCache
	local *cache.RankingLocalCache
}

func NewCachedRankingRepository(redis *cache.RankingRedisCache, local *cache.RankingLocalCache) RankingRepository {
	return &CachedRankingRepository{
		redis: redis,
		local: local,
	}
}

func (c *CachedRankingRepository) ReplaceTopN(ctx context.Context, arts []domain.Article) error {
	_ = c.local.Set(ctx, arts)
	return c.local.Set(ctx, arts)
}

func (c *CachedRankingRepository) GetTopN(ctx context.Context) ([]domain.Article, error) {
	data, err := c.local.Get(ctx)
	if err == nil {
		return data, nil
	}
	data, err = c.redis.Get(ctx)
	if err == nil {
		_ = c.local.Set(ctx, data)
	} else {
		return c.local.ForceGet(ctx)
	}
	return data, err
}
