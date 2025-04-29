package domain

import "time"

type Author struct {
	Id   int64
	Name string
}

type Article struct {
	Id       int64
	Title    string
	Content  string
	Author   Author
	Status   ArticleStatus
	Ctime    time.Time
	Utime    time.Time
	AuthorId int64
}

type ArticleStatus uint8

const (
	// ArticleStatusUnknown 为了避免零值之类的问题
	ArticleStatusUnknown ArticleStatus = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)
