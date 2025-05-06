package articles

import (
	"context"
	"github.com/zmsocc/practice/webook/internal/domain"
	"github.com/zmsocc/practice/webook/internal/repository/cache"
	"github.com/zmsocc/practice/webook/internal/repository/dao/articles"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, id, author int64, status domain.ArticleStatus) error
	GetByAuthor(ctx context.Context, uid int64, offset, limit int) ([]domain.Article, error)
}

type articleRepository struct {
	dao   articles.ArticleDAO
	cache cache.ArticleCache
}

func NewArticleRepository(dao articles.ArticleDAO) ArticleRepository {
	return &articleRepository{
		dao: dao,
	}
}

func (ar *articleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return ar.dao.Insert(ctx, articles.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	})
}

func (ar *articleRepository) Update(ctx context.Context, art domain.Article) error {
	return ar.dao.UpdateById(ctx, articles.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	})
}

func (ar *articleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	id, err := ar.dao.Sync(ctx, ar.toEntity(art))
	return id, err
}

func (ar *articleRepository) SyncStatus(ctx context.Context, id, author int64, status domain.ArticleStatus) error {
	return ar.dao.SyncStatus(ctx, id, author, status.ToUint8())
}

func (ar *articleRepository) GetByAuthor(ctx context.Context, uid int64, offset, limit int) ([]domain.Article, error) {
	panic("implement me")
}

func (ar *articleRepository) toEntity(art domain.Article) articles.Article {
	return articles.Article{
		Id:       art.Id,
		AuthorId: art.Author.Id,
		Title:    art.Title,
		Content:  art.Content,
		Status:   art.Status.ToUint8(),
		Ctime:    art.Ctime.UnixMilli(),
		Utime:    art.Utime.UnixMilli(),
	}
}

func (ar *articleRepository) toDomain(art articles.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Status:  domain.ArticleStatus(art.Status),
		Author: domain.Author{
			Id: art.AuthorId,
		},
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
}
