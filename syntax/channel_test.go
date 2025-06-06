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
		time.Sleep(1 * time.Second)
		ch1 <- 123
	}()
	go func() {
		time.Sleep(1 * time.Second)
		ch2 <- 456
	}()
	select {
	case val := <-ch1:
		t.Log("ch1", val)
		val = <-ch2
		t.Log("ch2", val)
	case val := <-ch2:
		t.Log("ch2", val)
		val = <-ch1
		t.Log("ch1", val)
	}
}

func TestLoopChannel(t *testing.T) {
	ch := make(chan int, 1)
	go func() {
		for i := 0; i < 3; i++ {
			ch <- i
			time.Sleep(1 * time.Second)
		}
		for i := 3; i < 6; i++ {
			ch <- i
			time.Sleep(1 * time.Second)
		}
		close(ch)
	}()
	//go func() {
	//	for i := 10; i < 13; i++ {
	//		ch <- i
	//		time.Sleep(1 * time.Second)
	//	}
	//	for i := 100; i < 103; i++ {
	//		ch <- i
	//		time.Sleep(1 * time.Second)
	//	}
	//	close(ch)
	//}()
	for val := range ch {
		t.Log(val)
	}
	t.Log("channel 被关了")
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
