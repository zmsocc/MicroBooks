package domain

import "time"

type Author struct {
	Id   int64
	Name string
}

type Article struct {
	Id      int64
	Title   string
	Content string
	// Author 要从用户来
	Author Author
	Status ArticleStatus
	Ctime  time.Time
	Utime  time.Time
	// 做成这样，就应该在 service 或者 repository 里面完成构造
	// 设计成这个样子，就认为 Interactive 是 Article 的一个属性（值对象）
	// Intr Interactive
}

type ArticleStatus uint8

const (
	// ArticleStatusUnknown 为了避免零值之类的问题
	ArticleStatusUnknown ArticleStatus = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

func (s ArticleStatus) ToUint8() uint8 {
	return uint8(s)
}
