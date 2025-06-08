package sarama

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var addrs = []string{"localhost:9094"}

func TestSyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(addrs, cfg)
	assert.NoError(t, err)
	p, offset, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: "read_article",
		// 消息数据本体
		// 转 JSON
		Value: sarama.StringEncoder(`{"aid": 1, "uid": 123}`),
		// 会在生产者和消费者之间传递
		//Headers: []sarama.RecordHeader{
		//	{
		//		Key:   []byte("trace_id"),
		//		Value: []byte("123456"),
		//	},
		//},
		//// 只作用于发送过程
		//Metadata: "这是 metadata",
	})
	assert.NoError(t, err)
	t.Log(p, offset)
}

func TestAsyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	producer, err := sarama.NewAsyncProducer(addrs, cfg)
	require.NoError(t, err)
	msgCh := producer.Input()
	go func() {
		for {
			msg := &sarama.ProducerMessage{
				Topic: "test_topic",
				Key:   sarama.StringEncoder("oid-123"),
				// 消息数据本体
				// 转 JSON
				Value: sarama.StringEncoder("hello, 这是一条消息 B"),
				// 会在生产者和消费者之间传递
				Headers: []sarama.RecordHeader{
					{
						Key:   []byte("trace_id"),
						Value: []byte("123456"),
					},
				},
				// 只作用于发送过程
				Metadata: "这是 metadata",
			}
			select {
			case msgCh <- msg:
				//default:
			}
		}
	}()
	errCh := producer.Errors()
	succCh := producer.Successes()
	for {
		select {
		case err := <-errCh:
			t.Log("发送出了问题", err.Err)
		case <-succCh:
			t.Log("发送成功")
		}
	}
}

type JSONEncoder struct {
	Data any
}
