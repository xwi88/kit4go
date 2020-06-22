package kafka

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/kdpujie/log4go"
)

func TestConsumerConsume(t *testing.T) {
	w := log4go.NewConsoleWriterWithLevel(log4go.DEBUG)
	w.SetColor(true)
	log4go.Register(w)

	config := cluster.NewConfig()
	config.Version = sarama.V2_5_0_0 // specify appropriate version
	config.Consumer.Return.Errors = true
	// config.Consumer.IsolationLevel = sarama.ReadCommitted // default ReadUncommitted=0
	// config.Consumer.Offsets.Initial = sarama.OffsetOldest // default OffsetOldest=-2
	// config.Consumer.Offsets.Initial = sarama.OffsetNewest // default OffsetOldest=-2
	config.Consumer.Group.Rebalance.Timeout = time.Second * 60 // default 60s
	// config.Consumer.Offsets.CommitInterval = time.Second // FIXME: must set
	config.Consumer.Offsets.CommitInterval = config.Consumer.Offsets.AutoCommit.Interval // FIXME: must set
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange               // default BalanceStrategyRange
	// config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin // default BalanceStrategyRange
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky // default BalanceStrategyRange
	config.Consumer.Group.Rebalance.Timeout = time.Second * 5
	config.Consumer.Group.Rebalance.Retry.Max = 4
	config.Consumer.Group.Rebalance.Retry.Backoff = time.Second * 2
	config.Consumer.Group.Heartbeat.Interval = time.Second * 6
	config.Consumer.Group.Session.Timeout = time.Second * 30
	config.Consumer.Offsets.AutoCommit.Enable = false

	brokers := []string{
		"10.14.41.57:9092",
		"10.14.41.58:9092",
		"10.14.41.59:9092",
	}
	topics := []string{"d-application-sys-log"}
	groupID := "dsp_will_20200601"

	group, err := NewConsumer(brokers, topics, groupID, config)
	if err != nil {
		panic(err)
	}
	defer func() { _ = group.Close() }() // error ...

	// Track errors
	go func() {
		for err := range group.c.Errors() {
			log4go.Info("[TestConsumerGroupConsume] consume errors:%v, type:%T, addr:%p", err, err, err)
		}
	}()

	// Iterate over consumer sessions.
	go group.StartConsumer(testConsumer)
	// time.Sleep(time.Second * 30)
	select {}

}

func testConsumer(msg *sarama.ConsumerMessage) {
	data := make(map[string]interface{})
	_ = json.Unmarshal(msg.Value, &data)
	log4go.Debug("[TestConsumerGroupConsume] testConsumer, topic:%v, partition:%v, offset:%v, timestamp:%v, "+
		"value:%+v",
		msg.Topic, msg.Partition, msg.Offset, msg.Timestamp, data)
	time.Sleep(time.Second * 1)
}
