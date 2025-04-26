package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/zmsocc/practice/webook/internal/domain"
	"github.com/zmsocc/practice/webook/internal/repository/cache"
	"github.com/zmsocc/practice/webook/internal/repository/dao"
	"time"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicate
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByID(ctx context.Context, id int64) (domain.User, error)
	Update(ctx context.Context, u domain.User) error
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
}

type userRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(d dao.UserDAO, c cache.UserCache) UserRepository {
	return &userRepository{
		dao:   d,
		cache: c,
	}
}

func (r *userRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
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
	user := r.entityToDomain(ue)
	err = r.cache.Set(ctx, user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *userRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *userRepository) Update(ctx context.Context, u domain.User) error {
	entity := r.domainToEntity(u)
	// 只更新指定字段
	return r.dao.Update(ctx, entity)
}

func (r *userRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Ctime:    u.Ctime.UnixMilli(),
	}
}

func (r *userRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}
