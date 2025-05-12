package test

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestChannelClose(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 123
	val, ok := <-ch
	t.Log(val, ok)
	close(ch)

	val, ok = <-ch
	t.Log(val, ok)
}

func TestForLoop(t *testing.T) {
	ch := make(chan int)
	go func() {
		for i := 0; i < 3; i++ {
			ch <- i
		}
		close(ch)
	}()
	for val := range ch {
		t.Log(val)
	}
	t.Log("发送完毕")
}

func TestSelect(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	go func() {
		for i := 0; i < 3; i++ {
			ch1 <- i
			ch2 <- i + 1
		}
		close(ch1)
		close(ch2)
	}()
	select {
	//case val := <-ch1:
	//	t.Log(val)
	case val := <-ch2:
		t.Log(val)
	}
}

func TestConsumer(t *testing.T) {
	cfg := sarama.NewConfig()
	eg, err := sarama.NewConsumerGroup([]string{"localhost:9094"}, "test_group", cfg)
	assert.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	err = eg.Consume(ctx, []string{"test_topic"}, &ConsumerHandler{})
	assert.NoError(t, err)
}

func TestProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{"localhost:9094"}, cfg)
	assert.NoError(t, err)
	p, offset, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder("hello world"),
		// 会在 producer 和 consumer 之间传递
		Headers: []sarama.RecordHeader{
			{Key: []byte("header1"), Value: []byte("header1_value")},
		},
		Metadata: map[string]any{"metadata1": "metadata_value1"},
	})
	assert.NoError(t, err)
	t.Log(p, offset)
}
