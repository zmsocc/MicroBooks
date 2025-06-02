package domain

import "time"

type Interactive struct {
	Id         int64
	Biz        string
	BizId      int64
	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64
	Utime      time.Time
	Ctime      time.Time
	Liked      bool
	Collected  bool
}
