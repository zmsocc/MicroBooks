package repository

import "context"

type HistoryRecordRepository interface {
	AddRecord(ctx context.Context, aid, uid int64) error
}
