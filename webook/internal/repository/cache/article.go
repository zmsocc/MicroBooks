package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/zmsocc/practice/webook/internal/domain"
	"time"
)

type ArticleCache interface {
	DelFirstPage(ctx context.Context, author int64) error
	DelPub(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (domain.Article, error)
	GetFirstPage(ctx context.Context, author int64) ([]domain.Article, error)
	GetPub(ctx context.Context, id int64) (domain.Article, error)
	Set(ctx context.Context, art domain.Article) error
	SetFirstPage(ctx context.Context, author int64, arts []domain.Article) error
	SetPub(ctx context.Context, art domain.Article) error
}

type RedisArticleCache struct {
	cmd redis.Cmdable
}

func NewArticleCache(cmd redis.Cmdable) ArticleCache {
	return &RedisArticleCache{
		cmd: cmd,
	}
}

func (c *RedisArticleCache) Get(ctx context.Context, id int64) (domain.Article, error) {
	data, err := c.cmd.Get(ctx, c.authorArtKey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var art domain.Article
	err = json.Unmarshal(data, &art)
	return art, err
}

func (c *RedisArticleCache) Set(ctx context.Context, art domain.Article) error {
	data, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, c.authorArtKey(art.Id), data, time.Minute).Err()
}

func (c *RedisArticleCache) GetFirstPage(ctx context.Context, author int64) ([]domain.Article, error) {
	data, err := c.cmd.Get(ctx, c.firstPageKey(author)).Bytes()
	if err != nil {
		return nil, err
	}
	var arts []domain.Article
	err = json.Unmarshal(data, &arts)
	return arts, err
}

func (c *RedisArticleCache) SetFirstPage(ctx context.Context, author int64, arts []domain.Article) error {
	for i := range arts {
		arts[i].Content = arts[i].Abstract()
	}
	data, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, c.firstPageKey(author), data, time.Minute*10).Err()
}

func (c *RedisArticleCache) DelFirstPage(ctx context.Context, author int64) error {
	return c.cmd.Del(ctx, c.firstPageKey(author)).Err()
}

func (c *RedisArticleCache) GetPub(ctx context.Context, id int64) (domain.Article, error) {
	data, err := c.cmd.Get(ctx, c.readerArtKey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var art domain.Article
	err = json.Unmarshal(data, &art)
	return art, err
}

func (c *RedisArticleCache) SetPub(ctx context.Context, art domain.Article) error {
	data, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, c.readerArtKey(art.Id), data, time.Minute*30).Err()
}

func (c *RedisArticleCache) DelPub(ctx context.Context, id int64) error {
	return c.cmd.Del(ctx, c.readerArtKey(id)).Err()
}

func (c *RedisArticleCache) authorArtKey(id int64) string {
	return fmt.Sprintf("article:author:%d", id)
}

func (c *RedisArticleCache) firstPageKey(id int64) string {
	return fmt.Sprintf("article:first_page:%d", id)
}

func (c *RedisArticleCache) readerArtKey(id int64) string {
	return fmt.Sprintf("article:reader:%d", id)
}
