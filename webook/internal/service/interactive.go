package service

import (
	"context"
	"github.com/zmsocc/practice/webook/internal/domain"
	"github.com/zmsocc/practice/webook/internal/repository"
	"golang.org/x/sync/errgroup"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, biz string, bizId, uid int64) error
	CancelLike(ctx context.Context, biz string, bizId, uid int64) error
	Collect(ctx context.Context, biz string, bizId, uid int64) error
	CancelCollect(ctx context.Context, biz string, bizId, uid int64) error
	Get(ctx context.Context, biz string, bizId, uid int64) (domain.Interactive, error)
}

type interactiveService struct {
	repo repository.InteractiveRepository
}

func NewInteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &interactiveService{
		repo: repo,
	}
}

func (i *interactiveService) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return i.repo.IncrReadCnt(ctx, biz, bizId)
}

func (i *interactiveService) Like(ctx context.Context, biz string, bizId, uid int64) error {
	return i.repo.IncrLike(ctx, biz, bizId, uid)
}

func (i *interactiveService) CancelLike(ctx context.Context, biz string, bizId, uid int64) error {
	return i.repo.DecrLike(ctx, biz, bizId, uid)
}

func (i *interactiveService) Collect(ctx context.Context, biz string, bizId, uid int64) error {
	return i.repo.AddCollectionItem(ctx, biz, bizId, uid)
}

func (i *interactiveService) CancelCollect(ctx context.Context, biz string, bizId, uid int64) error {
	return i.repo.DecrCollection(ctx, biz, bizId, uid)
}

func (i *interactiveService) Get(ctx context.Context, biz string, bizId, uid int64) (domain.Interactive, error) {
	var (
		eg        errgroup.Group
		intr      domain.Interactive
		liked     bool
		collected bool
	)
	eg.Go(func() error {
		var err error
		intr, err = i.repo.Get(ctx, biz, bizId)
		return err
	})
	eg.Go(func() error {
		var err error
		liked, err = i.repo.Liked(ctx, biz, bizId, uid)
		return err
	})
	eg.Go(func() error {
		var err error
		liked, err = i.repo.Collected(ctx, biz, bizId, uid)
		return err
	})
	err := eg.Wait()
	if err != nil {
		return domain.Interactive{}, err
	}
	intr.Liked = liked
	intr.Collected = collected
	return intr, err
}
