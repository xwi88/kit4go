package kafka

import (
	"context"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/kdpujie/log4go"
)

func TestConsumerConsume(t *testing.T) {
	w := log4go.NewConsoleWriterWithLevel(log4go.DEBUG)
	w.SetColor(true)
	log4go.Register(w)

	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0 // specify appropriate version
	config.Consumer.Return.Errors = true
	config.Consumer.IsolationLevel = sarama.ReadCommitted // default ReadUncommitted=0
	// config.Consumer.Offsets.Initial = sarama.OffsetOldest // default OffsetOldest=-2
	config.Consumer.Offsets.Initial = sarama.OffsetNewest // default OffsetOldest=-2
	// config.Consumer.Group.Rebalance.Timeout = time.Second*60 // default 60s
	// config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange      // default BalanceStrategyRange
	// config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin // default BalanceStrategyRange
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky // default BalanceStrategyRange
	config.Consumer.Group.Rebalance.Timeout = time.Second * 5
	config.Consumer.Group.Rebalance.Retry.Max = 4
	config.Consumer.Group.Rebalance.Retry.Backoff = time.Second * 2
	config.Consumer.Group.Heartbeat.Interval = time.Second * 6
	config.Consumer.Group.Session.Timeout = time.Second * 30

	brokers := []string{
		"10.14.41.57:9092",
		"10.14.41.58:9092",
		"10.14.41.59:9092",
	}
	topics := []string{"d-application-sys-log"}
	groupID := "dsp_will_20200601"

	group, err := NewConsumerGroup(brokers, topics, groupID, config)
	if err != nil {
		panic(err)
	}
	// defer func() { _ = group.Close() }() // error ...

	// Track errors
	// go func() {
	// 	for err := range group.cg.Errors() {
	// 		log4go.Info("[TestConsumerGroupConsume] consume errors:%v, type:%T, addr:%p", err, err, err)
	// 	}
	// }()

	// Iterate over consumer sessions.
	ctx := context.Background()

	// single run
	handler := exampleConsumerGroupHandler{}
	go group.StartConsumer(ctx, handler)
	time.Sleep(time.Second * 60)
	_ = group.Close()
	time.Sleep(time.Second * 5)

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
