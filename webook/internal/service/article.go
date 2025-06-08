package service

import (
	"context"
	"github.com/zmsocc/practice/webook/internal/domain"
	"github.com/zmsocc/practice/webook/internal/event/article"
	"github.com/zmsocc/practice/webook/internal/repository/articles"
	"github.com/zmsocc/practice/webook/pkg/logger"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, art domain.Article) error
	List(ctx context.Context, uid int64, offset, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPubById(ctx context.Context, id, uid int64) (domain.Article, error)
}

type articleService struct {
	repo     articles.ArticleRepository
	author   articles.ArticleAuthorRepository
	reader   articles.ArticleReaderRepository
	l        logger.Logger
	producer article.Producer
}

func NewArticleService(repo articles.ArticleRepository, l logger.Logger,
	producer article.Producer) ArticleService {
	return &articleService{
		repo:     repo,
		l:        l,
		producer: producer,
	}
}

func (svc *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnpublished
	if art.Id > 0 {
		err := svc.repo.Update(ctx, art)
		return art.Id, err
	}
	return svc.repo.Create(ctx, art)
}

func (svc *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	return svc.repo.Sync(ctx, art)
}

func (svc *articleService) Withdraw(ctx context.Context, art domain.Article) error {
	return svc.repo.SyncStatus(ctx, art.Id, art.Author.Id, domain.ArticleStatusPrivate)
}

func (svc *articleService) List(ctx context.Context, uid int64, offset, limit int) ([]domain.Article, error) {
	return svc.repo.List(ctx, uid, offset, limit)
}

func (svc *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return svc.repo.GetById(ctx, id)
}

func (svc *articleService) GetPubById(ctx context.Context, id, uid int64) (domain.Article, error) {
	_, err := svc.repo.GetPubById(ctx, id)
	if err == nil {
		go func() {
			er := svc.producer.ProduceReadEvent(
				ctx,
				article.ReadEvent{
					// 即便你的消费者要用 art 里面的数据，
					// 让他去查询，你不要在 event 里面带
					Uid: uid,
					Aid: id,
				})
			if er == nil {
				svc.l.Error("发送阅读者事件失败")
			}
		}()
	}
	return svc.repo.GetPubById(ctx, id)
}
