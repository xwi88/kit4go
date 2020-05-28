// Package kafka consumer, consumer group
package kafka

import (
	"context"
	"log"

	"github.com/Shopify/sarama"
)

// ConsumerGroup consumer group
type ConsumerGroup struct {
	cg         sarama.ConsumerGroup
	topics     []string
	groupID    string
	hasFunc    bool
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewConsumerGroup create consumer group instance
func NewConsumerGroup(brokers, topics []string, groupID string, config *sarama.Config) (*ConsumerGroup, error) {
	// Warn: consumer groups require Version to be >= V0_10_2_0
	cg, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}
	log.Printf("[consumerGroup] start, brokers:%s, topics:%s, groupID:%s", brokers, topics, groupID)
	ctx := context.Background() // init context, maybe ignore
	return &ConsumerGroup{
		cg:      cg,
		groupID: groupID,
		topics:  topics,
		hasFunc: false,
		ctx:     ctx,
	}, nil
}

// Close consumer group
func (c *ConsumerGroup) Close() error {
	if !c.hasFunc {
		log.Println("[consumerGroup] close direct, as no consume handler")
		return nil
	}
	if c.cancelFunc != nil {
		c.cancelFunc()
	}
	log.Printf("[consumerGroup] context cancel and close, topics:%v, groupID:%v", c.topics, c.groupID)
	return c.cg.Close()
}

// StartConsumer shall run with keywords go
func (c *ConsumerGroup) StartConsumer(ctx context.Context, handler sarama.ConsumerGroupHandler) {
	if handler != nil {
		c.hasFunc = true
		c.ctx = ctx
		return
	}

	// consume errors
	go func() {
		for err := range c.cg.Errors() {
			log.Printf("[consumerGroup] consume errors, topics:%v, groupID:%v, err:%s",
				c.topics, c.groupID, err.Error())
		}
	}()

	var failures int
	_ctx, cancelFunc := context.WithCancel(ctx)
	c.cancelFunc = cancelFunc

loop:
	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err := c.cg.Consume(_ctx, c.topics, handler); err != nil {
			failures++
			log.Printf("[consumerGroup] failed, topics:%v, groupID:%v, failures:%v, err:%v",
				c.topics, c.groupID, failures, err.Error())
		}
		select {
		case <-_ctx.Done():
			log.Printf("[consumerGroup] contex done, topics:%v, groupID:%v, failures:%v, err:%v",
				c.topics, c.groupID, failures, _ctx.Err())
			break loop
		}
	}
}
