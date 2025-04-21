package repository

import (
	"context"
	"errors"
	"github.com/zmsocc/practice/webook/internal/domain"
	"github.com/zmsocc/practice/webook/internal/repository/cache"
	"github.com/zmsocc/practice/webook/internal/repository/dao"
	"time"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(d *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   d,
		cache: c,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
	// 先从 cache 里面找, 再从 dao 里面找, 找到了回写 cache
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		return u, err
	}
	// 查看 Redis 是否限流了
	if ctx.Value("limited") == true {
		return domain.User{}, errors.New("触发限流，缓存未命中，不查询数据库")
	}
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	user := domain.User{
		Id:       ue.Id,
		Email:    ue.Email,
		Password: ue.Password,
	}
	err = r.cache.Set(ctx, user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *UserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		Ctime:    u.Ctime.UnixMilli(),
	}
}

func (r *UserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}
