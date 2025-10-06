package time

import (
	"context"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	// 用 ticker
	tm := time.NewTicker(time.Second)
	// 这一句不要忘了
	// 避免潜在的 goroutine 泄露问题
	defer tm.Stop()
	//for now := range tm.C {
	//	t.Log(now.Unix())
	//}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	for {
		// 方法2
	end:
		select {
		case <-ctx.Done():
			t.Log("超时了，或者被取消了")
			// break 不会推出循环
			//goto end 方法1
			break end // 方法2
		case now := <-tm.C:
			t.Log(now.Unix())
		}
	}
	// 方法1
	//end:
	//	t.Log("退出了循环")
}

func TestTime(t *testing.T) {
	tm := time.NewTimer(time.Second)
	defer tm.Stop()
	go func() {
		for now := range tm.C {
			t.Log(now.Unix())
		}
	}()
	time.Sleep(time.Second * 10)
}
