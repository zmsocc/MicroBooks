package gormx

import (
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"time"
)

type Callbacks struct {
	vector *prometheus2.SummaryVec
}

func NewCallbacks() *Callbacks {
	vector := prometheus2.NewSummaryVec(prometheus2.SummaryOpts{
		// 在这边，你要考虑设置各种 Namespace
		Namespace: "geekbang_daming",
		Subsystem: "webook",
		Name:      "gorm_query_time",
		Help:      "统计 GORM 的执行时间",
		ConstLabels: map[string]string{
			"db": "webook",
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.9:   0.01,
			0.99:  0.005,
			0.999: 0.0001,
		},
		// 如果是 JOIN 查询， table 就是 JOIN 在一起的
		// 或者 table 就是主表，A JOIN B， 记录的是 A
	}, []string{"type", "table"})
	pcb := &Callbacks{
		vector: vector,
	}
	prometheus2.MustRegister(vector)
	return pcb
}

func (c *Callbacks) RegisterAll(db *gorm.DB) {
	// 作用于 INSERT 语句
	err := db.Callback().Create().Before("*").
		Register("prometheus_create_before", c.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Create().After("*").
		Register("prometheus_create_after", c.after("create"))
	if err != nil {
		panic(err)
	}
	err = db.Callback().Update().Before("*").
		Register("prometheus_update_before", c.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Update().After("*").
		Register("prometheus_update_after", c.after("update"))
	if err != nil {
		panic(err)
	}
	err = db.Callback().Delete().Before("*").
		Register("prometheus_delete_before", c.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Delete().After("*").
		Register("prometheus_delete_after", c.after("delete"))
	if err != nil {
		panic(err)
	}
	err = db.Callback().Raw().Before("*").
		Register("prometheus_raw_before", c.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Raw().After("*").
		Register("prometheus_raw_after", c.after("raw"))
	if err != nil {
		panic(err)
	}
	err = db.Callback().Row().Before("*").
		Register("prometheus_row_before", c.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Row().After("*").
		Register("prometheus_row_after", c.after("row"))
	if err != nil {
		panic(err)
	}
}

func (c *Callbacks) before() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		startTime := time.Now()
		db.Set("start_time", startTime)
	}
}

func (c *Callbacks) after(typ string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		val, _ := db.Get("start_time")
		startTime, ok := val.(time.Time)
		if !ok {
			// 你啥都干不了
			return
		}
		table := db.Statement.Table
		if table == "" {
			table = "unknown"
		}
		c.vector.WithLabelValues(typ, table).
			Observe(float64(time.Since(startTime).Milliseconds()))
	}
}

func (c *Callbacks) Name() string {
	return "prometheus-query"
}

func (c *Callbacks) Initialize(db *gorm.DB) error {
	c.RegisterAll(db)
	return nil
}
