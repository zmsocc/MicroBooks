package ioc

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/zmsocc/practice/webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	type DBConfig struct {
		DSN string `yaml:"dsn"`
	}
	var cfg = DBConfig{
		DSN: "root:root@tcp(localhost:13336)/webook",
	}
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(fmt.Errorf("初始化配置失败: %s", err))
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
