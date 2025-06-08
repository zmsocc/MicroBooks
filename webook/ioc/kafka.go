package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"github.com/zmsocc/practice/webook/internal/event"
	"github.com/zmsocc/practice/webook/internal/event/article"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	var cfg = Config{
		Addrs: []string{"localhost:9094"},
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
func NewConsumers(c1 *article.InteractiveReadEventConsumer) []event.Consumer {
	return []event.Consumer{c1}
}
