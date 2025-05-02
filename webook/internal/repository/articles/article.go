package articles

import (
	"context"
	"github.com/zmsocc/practice/webook/internal/domain"
	"github.com/zmsocc/practice/webook/internal/repository/dao/articles"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
}

type articleRepository struct {
	dao articles.ArticleDAO
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
