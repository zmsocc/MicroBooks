package ioc

import (
	"github.com/zmsocc/practice/webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	type Config struct {
		DSN string `json:"dsn"`
	}
	var cfg = Config{
		DSN: "root:root@tcp(localhost:13336)/webook",
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
