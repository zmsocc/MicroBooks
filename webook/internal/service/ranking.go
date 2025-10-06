package service

import (
	"context"
	"errors"
	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	"github.com/zmsocc/practice/webook/internal/domain"
	"github.com/zmsocc/practice/webook/internal/repository"
	"math"
	"time"
)

type RankingService interface {
	TopN(ctx context.Context) error
	//TopN(ctx context.Context, n int64) error
	//TopN(ctx context.Context, n int64) ([]domain.Article, error)
}

type BatchRankingService struct {
	artSvc    ArticleService
	interSvc  InteractiveService
	repo      repository.RankingRepository
	batchSize int
	n         int
	// scoreFunc 不能返回负数
	scoreFunc func(t time.Time, likeCnt int64) float64
}

func NewBatchRankingService(artSvc ArticleService, interSvc InteractiveService) RankingService {
	return &BatchRankingService{
		artSvc:    artSvc,
		interSvc:  interSvc,
		batchSize: 100,
		n:         100,
		scoreFunc: func(t time.Time, likeCnt int64) float64 {
			ms := time.Since(t).Seconds()
			return float64(likeCnt-1) / math.Pow(ms+2, 1.5)
		},
	}
}

// TopN 准备分批
func (svc *BatchRankingService) TopN(ctx context.Context) error {
	arts, err := svc.topN(ctx)
	if err != nil {
		return err
	}
	// 在这里，存起来
	return svc.repo.ReplaceTopN(ctx, arts)
}

func (svc *BatchRankingService) topN(ctx context.Context) ([]domain.Article, error) {
	// 先拿一批数据
	now := time.Now()
	offset := 0
	type Score struct {
		art   domain.Article
		score float64
	}
	topN := queue.NewPriorityQueue(svc.n, func(src Score, dst Score) int {
		if src.score > dst.score {
			return 1
		} else if src.score == dst.score {
			return 0
		} else {
			return -1
		}
	})
	for {
		// 这里拿到了一批
		arts, err := svc.artSvc.ListPub(ctx, now, offset, svc.batchSize)
		if err != nil {
			return nil, err
		}
		ids := slice.Map[domain.Article, int64](arts,
			func(idx int, src domain.Article) int64 {
				return src.Id
			})
		// 要去找到对应的点赞数据
		inters, err := svc.interSvc.GetByIds(ctx, "article", ids)
		if err != nil {
			return nil, err
		}
		// 合并计算 score
		// 排序
		for _, art := range arts {
			inter := inters[art.Id]
			//if !ok {
			// 你都没有，肯定不可能是热榜
			//continue
			//}
			score := svc.scoreFunc(art.Utime, inter.LikeCnt)
			// 我要考虑，我这个 score 在不在前一百名
			// 拿到热度最低的
			err = topN.Enqueue(Score{
				art:   art,
				score: score,
			})
			// 小根堆已经满了
			if errors.Is(err, queue.ErrOutOfCapacity) {
				val, _ := topN.Dequeue()
				if val.score < score {
					err = topN.Enqueue(Score{
						art:   art,
						score: score,
					})
				} else {
					_ = topN.Enqueue(val)
				}
			}
		}
		// 一批已经处理完了，我要不要进入下一批?我怎么知道还有没有？
		if len(arts) < svc.batchSize || now.Sub(arts[len(arts)-1].Utime).Hours() > 24*7 {
			// 我这一批都没取够，就可以肯定没有下一批了
			// 或者已经取到七天之前的数据了，说明可以中断了
			break
		}
		// 这边要更新 offset
		offset += len(arts)
	}
	// 最后得出结果
	res := make([]domain.Article, svc.n)
	for i := svc.n - 1; i >= 0; i-- {
		val, err := topN.Dequeue()
		if err != nil {
			// 说明取完了，不够 n
			break
		}
		res[i] = val.art
	}
	return res, nil
}
