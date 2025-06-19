package saramax

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/zmsocc/practice/webook/pkg/logger"
	"time"
)

type BatchHandler[T any] struct {
	l  logger.Logger
	fn func(msgs []*sarama.ConsumerMessage, ts []T) error
	// 用 option 模式来设置这个 batchSize 和 duration
	batchSize     int
	batchDuration time.Duration
}

func NewBatchHandler[T any](l logger.Logger, fn func(msgs []*sarama.ConsumerMessage, ts []T) error) *BatchHandler[T] {
	return &BatchHandler[T]{
		fn:            fn,
		l:             l,
		batchSize:     10,
		batchDuration: time.Second,
	}
}

func (b *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgsCh := claim.Messages()
	batchSize := b.batchSize
	for {
		ctx, cancel := context.WithTimeout(context.Background(), b.batchDuration)
		done := false
		msgs := make([]*sarama.ConsumerMessage, 0, batchSize)
		ts := make([]T, 0, batchSize)
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				// 这边代表超时了
				done = true
			case msg, ok := <-msgsCh:
				if !ok {
					cancel()
					// 代表消费者关闭了
					return nil
				}
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					b.l.Error("反序列化失败", logger.Error(err),
						logger.String("topic", msg.Topic),
						logger.Int32("partition", msg.Partition),
						logger.Int64("offset", msg.Offset))
					continue
				}
				msgs = append(msgs, msg)
				ts = append(ts, t)
			}
		}
		cancel()
		if len(msgs) == 0 {
			continue
		}
		err := b.fn(msgs, ts)
		if err != nil {
			b.l.Error("调用业务接口失败", logger.Error(err))
		}
		// 就这样
		for _, msg := range msgs {
			// 这样万无一失
			session.MarkMessage(msg, "")
		}
	}
}
