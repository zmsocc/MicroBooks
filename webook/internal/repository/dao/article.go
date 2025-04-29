package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type ArticleDao struct {
	db *gorm.DB
}

func (d *ArticleDao) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := d.db.WithContext(ctx).Create(&art).Error
	// 返回自增主键
	return art.Id, err
}

func (d *ArticleDao) Update(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	res := d.db.Model(&Article{}).WithContext(ctx).
		Where("id = ? AND author_id = ?", art.Id, art.AuthorId).
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"utime":   now,
		})
	err := res.Error
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("更新数据失败")
	}
	return nil
}

type Article struct {
	Id       int64 `gorm:"primary_key;autoIncrement"`
	Title    string
	Content  string
	Ctime    int64
	Utime    int64
	AuthorId int64
}
