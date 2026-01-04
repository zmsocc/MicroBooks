package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"github.com/zmsocc/practice/webook/interactive/events"
	"github.com/zmsocc/practice/webook/internal/event"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	var cfg = Config{
		Addrs: []string{"127.0.0.1:9092"},
	}
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	client, err := sarama.NewClient(cfg.Addrs, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

func NewSyncProducer(client sarama.Client) sarama.SyncProducer {
	res, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return res
}

// NewConsumers 面临的问题依旧是所有的 Consumer 在这里注册一下
func NewConsumers(c1 *events.InteractiveReadEventBatchConsumer) []event.Consumer {
	return []event.Consumer{c1}
}
