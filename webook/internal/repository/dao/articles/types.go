package articles

import (
	"context"
	"errors"
)

var ErrPossibleIncorrectAuthor = errors.New("用户在尝试操作非本人数据")

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
	SyncStatus(ctx context.Context, id, author int64, status uint8) error
	FindByAuthor(ctx context.Context, uid int64, offset, limit int) ([]Article, error)
	GetById(ctx context.Context, id int64) (Article, error)
	GetPubById(ctx context.Context, id int64) (Article, error)
}
