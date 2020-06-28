// Package kafka consumer, normal consumer
package kafka

import (
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/kdpujie/log4go"
)

// Consumer simple consumer
type Consumer struct {
	c          *cluster.Consumer
	topics     []string
	groupID    string
	hasFunc    bool
	closeStart chan struct{}
	closeEnd   chan struct{}
}

// NewConsumer create consumer instance
func NewConsumer(brokers, topics []string, groupID string, config *cluster.Config) (*Consumer, error) {
	// init consumer
	consumer, err := cluster.NewConsumer(brokers, groupID, topics, config)
	if err != nil {
		return nil, err
	}
	log4go.Debug("[consumer] created, brokers:%s, topics:%s, groupID:%v", brokers, topics, groupID)
	return &Consumer{c: consumer, topics: topics, groupID: groupID, hasFunc: false,
		closeStart: make(chan struct{}), closeEnd: make(chan struct{}),
	}, nil
}

// Close consumer
func (c *Consumer) Close() error {
	if !c.hasFunc {
		log4go.Info("[consumer] close direct, as no consume func")
		return nil
	}
	c.closeStart <- struct{}{}
	<-c.closeEnd
	return c.c.Close()
}

// StartConsumer shall run with keywords go
func (c *Consumer) StartConsumer(fn func(*sarama.ConsumerMessage) error) {
	if fn != nil {
		c.hasFunc = true
	} else {
		log4go.Error("[consumer] consume failed, handler func nil, topics:%v, groupID:%v",
			c.topics, c.groupID)
		// avoid high frequency output, if in infinite loop
		time.Sleep(time.Second * 1)
		return
	}

	// consume errors
	go func() {
		for err := range c.c.Errors() {
			log4go.Error("[consumer] consume errors, topics:%v, groupID:%v, err:%v",
				c.topics, c.groupID, err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range c.c.Notifications() {
			log4go.Debug("[consumer] consume notifications, topics:%v, groupID:%v, notification:%+v",
				c.topics, c.groupID, ntf)
		}
	}()

	var failures int

loop:
	for {
		select {
		case msg, ok := <-c.c.Messages():
			if ok {
				if err := fn(msg); err == nil {
					// mark message as processed
					c.c.MarkOffset(msg, "")
				} else {
					failures++
					log4go.Error("[consumer] fn errors, count:%v, topics:%v, groupID:%v, err:%v",
						failures, c.topics, c.groupID, err.Error())
				}
			}
		case <-c.closeStart:
			log4go.Warn("[consumer] close, topics:%v, groupID:%v, failures:%v", c.topics, c.groupID, failures)
			break loop
		}
	}
	c.closeEnd <- struct{}{}
	log4go.Info("[consumer] close success, topics:%v, groupID:%v", c.topics, c.groupID)
}
