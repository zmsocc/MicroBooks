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
	//DelFirstPage(ctx context.Context, author int64) error
	Get(ctx context.Context, id int64) (domain.Article, error)
	//GetFirstPage(ctx context.Context, author int64) ([]domain.Article, error)
	Set(ctx context.Context, art domain.Article) error
	//SetFirstPage(ctx context.Context, author int64, arts []domain.Article) error
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

//func (c *RedisArticleCache) GetFirstPage(ctx context.Context, author int64) ([]domain.Article, error) {
//
//}
//
//func (c *RedisArticleCache) SetFirstPage(ctx context.Context, author int64, arts []domain.Article) error {
//
//}
//
//func (c *RedisArticleCache) DelFirstPage(ctx context.Context, author int64) error {
//
//}

func (c *RedisArticleCache) authorArtKey(id int64) string {
	return fmt.Sprintf("article:author:%d", id)
}
