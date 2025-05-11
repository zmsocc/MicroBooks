package repository

import (
	"context"
	"github.com/zmsocc/practice/webook/internal/repository/cache"
	"github.com/zmsocc/practice/webook/internal/repository/dao"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	IncrLike(ctx context.Context, biz string, bizId, uid int64) error
	DecrLike(ctx context.Context, biz string, bizId, uid int64) error
}

type interactiveRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
}

func NewInteractiveRepository(dao dao.InteractiveDAO, cache cache.InteractiveCache) InteractiveRepository {
	return &interactiveRepository{
		dao:   dao,
		cache: cache,
	}
}

func (i *interactiveRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	err := i.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	return i.cache.IncrReadCntIfPresent(ctx, biz, bizId)
}

func (i *interactiveRepository) IncrLike(ctx context.Context, biz string, bizId, uid int64) error {
	err := i.dao.IncrLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return i.cache.IncrLikeCntIfPresent(ctx, biz, bizId)
}

func (i *interactiveRepository) DecrLike(ctx context.Context, biz string, bizId, uid int64) error {
	err := i.dao.DecrLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return i.cache.DecrLikeCntIfPresent(ctx, biz, bizId)
}
