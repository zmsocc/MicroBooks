package test

import (
	"github.com/IBM/sarama"
	"log"
)

type ConsumerHandler struct{}

func NewConsumerHandler() *ConsumerHandler {
	return &ConsumerHandler{}
}

func (c *ConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	// 执行一些初始化的事情
	log.Println("Handler Setup")
	// 假设要重置到 0
	var offset int64 = 0
	// 遍历所有的分区
	partitions := session.Claims()["test_topic"]
	for _, p := range partitions {
		session.ResetOffset("test_topic", p, offset, "metadata")
	}
	return nil
}

func (c *ConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	// 执行一些清理工作
	log.Println("Handler Cleanup")
	return nil
}

func (c *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ch := claim.Messages()
	for msg := range ch {
		log.Println(msg)
		// 标记为消费成功
		session.MarkMessage(msg, "")
	}
	return nil
}
