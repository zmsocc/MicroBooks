package ioc

import (
	"github.com/zmsocc/practice/webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	var cfg = WebookConfig{
		DB: DBConfig{
			DSN: "root:root@tcp(localhost:13336)/webook",
		},
		Redis: RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       1,
		},
	}
	db, err := gorm.Open(mysql.Open(cfg.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

type DBConfig struct {
	DSN string `json:"dsn"`
}

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type WebookConfig struct {
	DB    DBConfig
	Redis RedisConfig
}
