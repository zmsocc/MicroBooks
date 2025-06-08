package sarama

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"log"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	cfg := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(addrs, "test_consumer", cfg)
	require.NoError(t, err)

	// 带超时的 context
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = consumer.Consume(ctx, []string{"read_article"}, testConsumerGroupHandler{})
	// 你消费结束，就会到这里
	t.Log(err, time.Since(start).String())
}

type testConsumerGroupHandler struct {
}

func (t testConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	partitions := session.Claims()["test_topic"]
	for _, part := range partitions {
		session.ResetOffset("test_topic", part, sarama.OffsetOldest, "")
	}
	return nil
}

func (t testConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Println("Cleanup session:")
	return nil
}

// 代表的是你和 Kafka 的会话(从建立连接到连接彻底断掉的那一段时间)
func (t testConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	//for msg := range msgs {
	//	//var bizMsg MyBizMsg
	//	//err := json.Unmarshal(msg.Value, &bizMsg)
	//	//if err != nil {
	//	//	// 这就是消费信息出错
	//	//	// 大多数时候就是重试
	//	//	// 记录日志
	//	//	continue
	//	//}
	//	m1 := msg
	//	go func() {
	//		// 消费 msg
	//		log.Println(string(m1.Value))
	//		session.MarkMessage(m1, "")
	//	}()
	//}
	const batchSize = 10
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		var eg errgroup.Group
		var last *sarama.ConsumerMessage
		done := false
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				// 这边代表超时了
				done = true
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					// 代表消费者关闭了
					return nil
				}
				last = msg
				eg.Go(func() error {
					// 我就在这里消费
					time.Sleep(time.Second)
					log.Println(string(msg.Value))
					return nil
				})
			}
		}
		cancel()
		err := eg.Wait()
		if err != nil {
			// 这边能怎么办
			// 记录日志
			continue
		}
		// 就这样
		if last != nil {
			session.MarkMessage(last, "")
		}
	}
	// 什么情况下会到这里
	// msgs 被人关了
}

//type MyBizMsg struct {
//	Name string
//}
//
//// 返回只读的 channel
//// 优先用这个
//func ChannelV1() <-chan struct{} {
//	panic("implement me")
//}
//
//// 返回可读可写的 channel
//func ChannelV2() <-chan struct{} {
//	panic("implement me")
//}
//
//// 返回只写 channel
//func ChannelV3() <-chan struct{} {
//	panic("implement me")
//}
