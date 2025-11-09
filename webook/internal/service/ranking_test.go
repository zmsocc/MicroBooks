package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zmsocc/practice/webook/internal/domain"
	svcmocks "github.com/zmsocc/practice/webook/internal/service/mocks"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestRankingTopN(t *testing.T) {
	now := time.Now()
	testCase := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (ArticleService, InteractiveService)
		wantErr  error
		wantArts []domain.Article
	}{
		{
			name: "计算成功",
			// 怎么模拟我的数据？
			mock: func(ctrl *gomock.Controller) (ArticleService, InteractiveService) {
				artSvc := svcmocks.NewMockArticleService(ctrl)
				// 最简单，一批就搞完
				artSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 0, 3).
					Return([]domain.Article{
						{Id: 1, Utime: now, Ctime: now},
						{Id: 2, Utime: now, Ctime: now},
						{Id: 3, Utime: now, Ctime: now},
					}, nil)
				artSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 3, 3).
					Return([]domain.Article{}, nil)
				interSvc := svcmocks.NewMockInteractiveService(ctrl)
				interSvc.EXPECT().GetByIds(gomock.Any(), "article", []int64{1, 2, 3}).
					Return(map[int64]domain.Interactive{
						1: {BizId: 1, LikeCnt: 1},
						2: {BizId: 2, LikeCnt: 2},
						3: {BizId: 3, LikeCnt: 3},
					}, nil)
				interSvc.EXPECT().GetByIds(gomock.Any(), "article", []int64{}).
					Return(map[int64]domain.Interactive{}, nil)
				return artSvc, interSvc
			},
			wantArts: []domain.Article{
				{Id: 3, Utime: now, Ctime: now},
				{Id: 2, Utime: now, Ctime: now},
				{Id: 1, Utime: now, Ctime: now},
			},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			artSvc, interSvc := tc.mock(ctrl)
			svc := NewBatchRankingService(artSvc, interSvc).(*BatchRankingService)
			// 为了测试
			svc.batchSize = 3
			svc.n = 3
			svc.scoreFunc = func(t time.Time, likeCnt int64) float64 {
				return float64(likeCnt)
			}
			arts, err := svc.topN(context.Background())
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantArts, arts)
		})
	}
}
