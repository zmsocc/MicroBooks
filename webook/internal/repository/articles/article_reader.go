package articles

import (
	"context"
	"github.com/zmsocc/practice/webook/internal/domain"
)

type ArticleReaderRepository interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
}
