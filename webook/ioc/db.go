package ioc

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/zmsocc/practice/webook/internal/repository/dao"
	"github.com/zmsocc/practice/webook/pkg/gormx"
	"github.com/zmsocc/practice/webook/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"gorm.io/plugin/prometheus"
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
	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "webook",
		RefreshInterval: 15,
		StartServer:     false,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"thread_running"},
			},
		},
	}))
	if err != nil {
		panic(err)
	}
	// 监控查询的执行时间
	pcb := gormx.NewCallbacks()
	//pcb.RegisterAll(db)
	err = db.Use(pcb)
	if err != nil {
		return nil
	}
	db.Use(tracing.NewPlugin(tracing.WithDBName("webook"),
		tracing.WithQueryFormatter(func(query string) string {
			l := logger.NewNopLogger()
			l.Debug("", logger.String("query", query))
			return query
		}),
		// 不要记录 metrics
		tracing.WithoutMetrics(),
		// 不要记录查询参数
		tracing.WithoutQueryVariables()))

	return db
}
