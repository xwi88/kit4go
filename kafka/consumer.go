// Package kafka consumer, normal consumer
package kafka

import (
	"log"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
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
	log.Printf("[consumer] start, brokers:%s, topics:%s, groupID:%v", brokers, topics, groupID)
	return &Consumer{c: consumer, topics: topics, groupID: groupID, hasFunc: false,
		closeStart: make(chan struct{}), closeEnd: make(chan struct{}),
	}, nil
}

// Close consumer
func (c *Consumer) Close() error {
	if !c.hasFunc {
		log.Println("[consumer] close direct, as no consume func")
		return nil
	}
	c.closeStart <- struct{}{}
	<-c.closeEnd
	return c.c.Close()
}

// StartConsumer shall run with keywords go
func (c *Consumer) StartConsumer(fn func(*sarama.ConsumerMessage)) {
	if fn != nil {
		c.hasFunc = true
	}

	// consume errors
	go func() {
		for err := range c.c.Errors() {
			log.Printf("[consumer] consume errors, topics:%v, groupID:%v, err:%v",
				c.topics, c.groupID, err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range c.c.Notifications() {
			log.Printf("[consumer] consume notifications, topics:%v, groupID:%v, notification:%v",
				c.topics, c.groupID, ntf)
		}
	}()

	var failures int

loop:
	for {
		select {
		case msg, ok := <-c.c.Messages():
			if ok {
				fn(msg)
				// mark message as processed
				c.c.MarkOffset(msg, "")
				failures++
			}
		case <-c.closeStart:
			break loop
		}
	}
	log.Printf("[consumer] failed, topics:%v, groupID:%v, failures:%v",
		c.topics, c.groupID, failures)
	c.closeEnd <- struct{}{}
}
