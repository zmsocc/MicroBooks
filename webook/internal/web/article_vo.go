package web

import (
	"github.com/zmsocc/practice/webook/internal/domain"
)

type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (req ArticleReq) toDomain(uid int64) domain.Article {
	return domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uid,
		},
	}
}

type ListReq struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type LikeReq struct {
	Id int64 `json:"id"`
	// 点赞，取消点赞都复用这个
	Like bool `json:"like"`
}

type ArticleVO struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Abstract string `json:"abstract"`
	// Author 要从用户来
	Author string `json:"author"`
	Status uint8  `json:"status"`
	Ctime  string `json:"ctime"`
	Utime  string `json:"utime"`

	// 点赞之类的信息
	ReadCnt    int64 `json:"read_cnt"`
	CollectCnt int64 `json:"collect_cnt"`
	LikeCnt    int64 `json:"like_cnt"`

	// 个人是否点赞信息
	Collected bool `json:"collected"`
	Liked     bool `json:"liked"`
}
