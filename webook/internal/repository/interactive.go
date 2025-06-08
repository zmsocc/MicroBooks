package repository

import (
	"context"
	"errors"
	"github.com/zmsocc/practice/webook/internal/domain"
	"github.com/zmsocc/practice/webook/internal/repository/cache"
	"github.com/zmsocc/practice/webook/internal/repository/dao"
	"github.com/zmsocc/practice/webook/pkg/logger"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	BatchIncrReadCnt(ctx context.Context, biz []string, bizId []int64) error
	IncrLike(ctx context.Context, biz string, bizId, uid int64) error
	DecrLike(ctx context.Context, biz string, bizId, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, bizId, uid int64) error
	DecrCollection(ctx context.Context, biz string, bizId, uid int64) error
	Liked(ctx context.Context, biz string, bizId, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, bizId, uid int64) (bool, error)
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
}

type interactiveRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
	l     logger.Logger
}

func NewInteractiveRepository(dao dao.InteractiveDAO, cache cache.InteractiveCache,
	l logger.Logger) InteractiveRepository {
	return &interactiveRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

func (i *interactiveRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	err := i.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	return i.cache.IncrReadCntIfPresent(ctx, biz, bizId)
}

// BatchIncrReadCnt bizs 和 ids 的长度必须相等
func (i *interactiveRepository) BatchIncrReadCnt(ctx context.Context, biz []string, bizId []int64) error {
	// 在这里要不要检测 bizs 和 ids 的长度是否相等
	return i.dao.BatchIncrReadCnt(ctx, biz, bizId)
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

func (i *interactiveRepository) AddCollectionItem(ctx context.Context, biz string, bizId, uid int64) error {
	err := i.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
		Biz:   biz,
		BizId: bizId,
		Uid:   uid,
	})
	if err != nil {
		return err
	}
	return i.cache.IncrCollectCntIfPresent(ctx, biz, bizId)
}

func (i *interactiveRepository) DecrCollection(ctx context.Context, biz string, bizId, uid int64) error {
	err := i.dao.DecrCollectInfo(ctx, dao.UserCollectionBiz{
		Biz:   biz,
		BizId: bizId,
		Uid:   uid,
	})
	if err != nil {
		return err
	}
	return i.cache.DecrCollectCntIfPresent(ctx, biz, bizId)
}

func (i *interactiveRepository) Liked(ctx context.Context, biz string, bizId, uid int64) (bool, error) {
	_, err := i.dao.GetLikeInfo(ctx, biz, bizId, uid)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, dao.ErrRecordNotFound):
		return false, nil
	default:
		return false, err
	}
}

func (i *interactiveRepository) Collected(ctx context.Context, biz string, bizId, uid int64) (bool, error) {
	_, err := i.dao.GetCollectInfo(ctx, biz, bizId, uid)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, dao.ErrRecordNotFound):
		return false, nil
	default:
		return false, err
	}
}

func (i *interactiveRepository) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	intr, err := i.cache.Get(ctx, biz, bizId)
	if err == nil {
		return intr, nil
	}
	// 在这里查询数据库
	daoIntr, err := i.dao.Get(ctx, biz, bizId)
	if err != nil {
		return domain.Interactive{}, err
	}
	intr = i.toDomain(daoIntr)
	go func() {
		er := i.cache.Set(ctx, biz, bizId, intr)
		// 记录日志
		if er != nil {
			i.l.Error("会写缓存失败",
				logger.String("biz", biz),
				logger.Int64("bizId", bizId))
		}
	}()
	return intr, nil
}

func (i *interactiveRepository) toDomain(dao dao.Interactive) domain.Interactive {
	return domain.Interactive{
		ReadCnt:    dao.ReadCnt,
		LikeCnt:    dao.LikeCnt,
		CollectCnt: dao.CollectCnt,
	}
}
