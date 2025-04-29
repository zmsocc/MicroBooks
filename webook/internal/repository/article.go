package repository

import (
	"context"
	"github.com/zmsocc/practice/webook/internal/domain"
	"github.com/zmsocc/practice/webook/internal/repository/dao"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
}

type articleRepository struct {
	dao *dao.ArticleDao
}

func NewArticleRepository(dao *dao.ArticleDao) ArticleRepository {
	return &articleRepository{
		dao: dao,
	}
}

func (ar *articleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return ar.dao.Insert(ctx, dao.Article{
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.AuthorId,
		Status:   uint8(art.Status),
	})
}

func (ar *articleRepository) Update(ctx context.Context, art domain.Article) error {
	return ar.dao.Update(ctx, ar.toEntity(art))
}

func (ar *articleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		AuthorId: art.AuthorId,
		Title:    art.Title,
		Content:  art.Content,
		Ctime:    art.Ctime.UnixMilli(),
		Utime:    art.Utime.UnixMilli(),
	}
}

func (ar *articleRepository) toDomain(art dao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Ctime:   time.UnixMilli(art.Ctime),
		Utime:   time.UnixMilli(art.Utime),
	}
}
