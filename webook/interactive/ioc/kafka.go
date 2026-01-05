package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"github.com/zmsocc/practice/webook/interactive/events"
	"github.com/zmsocc/practice/webook/pkg/saramax"
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

func NewConsumers(c1 *events.InteractiveReadEventConsumer) []saramax.Consumer {
	return []saramax.Consumer{c1}
}
