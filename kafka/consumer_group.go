// Package kafka consumer, consumer group
package kafka

import (
	"context"
	"time"

	"github.com/Shopify/sarama"
	"github.com/kdpujie/log4go"
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
	log4go.Debug("[consumerGroup] created, brokers:%s, topics:%s, groupID:%s", brokers, topics, groupID)
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
		log4go.Info("[consumerGroup] close direct, as no consume handler")
		return nil
	}
	if c.cancelFunc != nil {
		c.cancelFunc()
	}
	log4go.Info("[consumerGroup] close called, topics:%v, groupID:%v, cancelFunc:%v",
		c.topics, c.groupID, c.cancelFunc)
	return nil
}

// StartConsumer shall run with keywords go
func (c *ConsumerGroup) StartConsumer(ctx context.Context, handler sarama.ConsumerGroupHandler) {
	if handler != nil {
		c.hasFunc = true
		c.ctx = ctx
	} else {
		log4go.Error("[consumerGroup] start consume failed, handler nil, topics:%v, groupID:%v",
			c.topics, c.groupID)
		// avoid high frequency output, if in infinite loop
		time.Sleep(time.Second * 1)
		return
	}

	// consume errors
	go func() {
		for err := range c.cg.Errors() {
			log4go.Error("[consumerGroup] consume errors, topics:%v, groupID:%v, err:%s",
				c.topics, c.groupID, err.Error())
		}
	}()

	var failures int
	_ctx, cancelFunc := context.WithCancel(ctx)
	c.cancelFunc = cancelFunc
	log4go.Debug("[consumerGroup] bind cancelFun, topics:%v, groupID:%v, cancelFun:%v",
		c.topics, c.groupID, cancelFunc)

loop:
	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err := c.cg.Consume(_ctx, c.topics, handler); err != nil {
			failures++
			log4go.Error("[consumerGroup] consume failed, topics:%v, groupID:%v, failures:%v, err:%v",
				c.topics, c.groupID, failures, err.Error())
			if err == sarama.ErrClosedConsumerGroup {
				// consumer group chan closed
				log4go.Error("[consumerGroup] consume error, topics:%v, groupID:%v, failures:%v, err:%v",
					c.topics, c.groupID, failures, "closed consumer group")
				time.Sleep(time.Second) // avoid frequency output in infinite loop!
				break loop
			}
		} else {
			log4go.Warn("[consumerGroup] consume exist, topics:%v, groupID:%v",
				c.topics, c.groupID)
		}

		select {
		case <-_ctx.Done():
			log4go.Info("[consumerGroup] context done, topics:%v, groupID:%v, failures:%v, err:%v",
				c.topics, c.groupID, failures, _ctx.Err())
			break loop
		}
	}
	if err := c.cg.Close(); err != nil {
		log4go.Error("[consumerGroup] close failed, topics:%v, groupID:%v, failures:%v, err:%v",
			c.topics, c.groupID, failures, err.Error())
	} else {
		log4go.Info("[consumerGroup] close success, topics:%v, groupID:%v", c.topics, c.groupID)
	}
}
