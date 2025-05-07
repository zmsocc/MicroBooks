package articles

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/zmsocc/practice/webook/internal/domain"
	"github.com/zmsocc/practice/webook/internal/repository"
	"github.com/zmsocc/practice/webook/internal/repository/cache"
	"github.com/zmsocc/practice/webook/internal/repository/dao/articles"
	"github.com/zmsocc/practice/webook/pkg/logger"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, id, author int64, status domain.ArticleStatus) error
	List(ctx context.Context, uid int64, offset, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPubById(ctx context.Context, id int64) (domain.Article, error)
}

type articleRepository struct {
	dao      articles.ArticleDAO
	artCache cache.ArticleCache
	l        logger.Logger
	userRepo repository.UserRepository
}

func NewArticleRepository(dao articles.ArticleDAO, artCache cache.ArticleCache,
	l logger.Logger) ArticleRepository {
	return &articleRepository{
		dao:      dao,
		artCache: artCache,
		l:        l,
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
	if err == nil {
		er := ar.artCache.DelFirstPage(ctx, art.Id)
		if er != nil {
			return 0, er
		}
		e := ar.artCache.SetPub(ctx, art)
		if e != nil {
			// 不需要特别关心
			ar.l.Warn("设置缓存失败")
		}
	}
	return id, err
}

func (ar *articleRepository) SyncStatus(ctx context.Context, id, author int64, status domain.ArticleStatus) error {
	return ar.dao.SyncStatus(ctx, id, author, status.ToUint8())
}

func (ar *articleRepository) List(ctx context.Context, uid int64, offset, limit int) ([]domain.Article, error) {
	if offset == 0 && limit == 100 {
		data, err := ar.artCache.GetFirstPage(ctx, uid)
		if err == nil {
			go func() {
				ar.preCache(ctx, data)
			}()
			return data, err
		}
	}
	res, err := ar.dao.FindByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}
	data := slice.Map[articles.Article, domain.Article](res, func(idx int, src articles.Article) domain.Article {
		return ar.toDomain(src)
	})
	go func() {
		err = ar.artCache.SetFirstPage(ctx, uid, data)
		ar.l.Error("回写缓存失败", logger.Error(err))
		ar.preCache(ctx, data)
	}()
	return data, nil
}

func (ar *articleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {
	art, err := ar.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return ar.toDomain(art), nil
}

func (ar *articleRepository) GetPubById(ctx context.Context, id int64) (domain.Article, error) {
	art, err := ar.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	user, err := ar.userRepo.FindByID(ctx, art.AuthorId)
	if err != nil {
		return domain.Article{}, err
	}
	res := domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Status:  domain.ArticleStatus(art.Status),
		Author: domain.Author{
			Id:   user.Id,
			Name: user.Nickname,
		},
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
	return res, nil
}

func (ar *articleRepository) preCache(ctx context.Context, arts []domain.Article) {
	if len(arts) > 0 && len(arts[0].Content) < 1024*1024 {
		err := ar.artCache.Set(ctx, arts[0])
		if err != nil {
			ar.l.Error("提前预加载缓存失败", logger.Error(err))
		}
	}
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
