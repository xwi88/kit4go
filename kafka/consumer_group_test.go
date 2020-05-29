package kafka

import (
	"context"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/kdpujie/log4go"
)

type exampleConsumerGroupHandler struct{}

func (exampleConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	log4go.Info("[exampleConsumerGroupHandler] Setup")
	return nil
}
func (exampleConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	log4go.Info("[exampleConsumerGroupHandler] Cleanup")
	return nil
}
func (h exampleConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log4go.Info("[exampleConsumerGroupHandler] topic:%q partition:%d offset:%d", msg.Topic, msg.Partition, msg.Offset)
		sess.MarkMessage(msg, "")
		time.Sleep(time.Second * 5)
	}
	return nil
}

func TestConsumerGroupConsume(t *testing.T) {
	w := log4go.NewConsoleWriterWithLevel(log4go.DEBUG)
	w.SetColor(true)
	log4go.Register(w)

	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0 // specify appropriate version
	config.Consumer.Return.Errors = true

	brokers := []string{
		"10.14.41.57:9092",
		"10.14.41.58:9092",
		"10.14.41.59:9092",
	}
	topics := []string{"dsp_application_log"}
	groupID := "dsp_will_20200601"

	group, err := NewConsumerGroup(brokers, topics, groupID, config)
	if err != nil {
		panic(err)
	}
	defer func() { _ = group.Close() }()

	// Track errors
	go func() {
		for err := range group.cg.Errors() {
			log4go.Info("[TestConsumerGroupConsume] consume errors:%v, type:%T, addr:%p", err, err, err)
		}
	}()

	go func() {
		tk := time.NewTimer(time.Second * 3)
		defer tk.Stop()
		select {
		case <-tk.C:
			if err := group.Close(); err != nil {
				log4go.Error("[TestConsumerGroupConsume] timer close, err:%v", err.Error())
			}
			break
		}
	}()
	// Iterate over consumer sessions.
	ctx := context.Background()

	// single run
	handler := exampleConsumerGroupHandler{}
	group.StartConsumer(ctx, handler)

	// forbidden infinite loop called
	// for {
	// 	handler := exampleConsumerGroupHandler{}
	// 	// `Consume` should be called inside an infinite loop, when a
	// 	// server-side rebalance happens, the consumer session will need to be
	// 	// recreated to get the new claims
	// 	group.StartConsumer(ctx, handler)
	// 	// group.StartConsumer(ctx, nil) // also deal for test, but you shall not input the nil
	// }
}
