package repository

import (
	"context"
	"github.com/zmsocc/practice/webook/internal/repository/cache"
)

var (
	ErrCodeSendTooMany   = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany
)

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type codeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(c cache.CodeCache) CodeRepository {
	return &codeRepository{
		cache: c,
	}
}

func (repo *codeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}

func (repo *codeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return repo.cache.Get(ctx, biz, phone, code)
}
