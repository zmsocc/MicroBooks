package dao

import (
	"github.com/zmsocc/practice/webook/internal/repository/dao/articles"
	"gorm.io/gorm"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &articles.Article{})
}
