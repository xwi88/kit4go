// Package kafka producer, sync producer
package kafka

import (
	"log"

	"github.com/Shopify/sarama"

	"github.com/xwi88/kit4go/utils"
)

var (
	// DefaultSyncProducerBufferSize default sync buffer size
	DefaultSyncProducerBufferSize = 1
)

// SyncProducer sync producer
type SyncProducer struct {
	producer sarama.SyncProducer
	messages chan *sarama.ProducerMessage
	cfg      *sarama.Config
	stop     chan struct{}
}

// NewSyncProducer create sync producer instance
func NewSyncProducer(brokers []string, bufferSize int, cfg *sarama.Config) (sp *SyncProducer, err error) {
	sp = new(SyncProducer)
	// if not set, use kafka default ChannelBufferSize=256
	if bufferSize == 0 {
		bufferSize = cfg.ChannelBufferSize
	} else if bufferSize <= DefaultSyncProducerBufferSize {
		// if set, but not greater than our limit DefaultAsyncProducerBufferSize, set it
		bufferSize = DefaultSyncProducerBufferSize
	}
	sp.messages = make(chan *sarama.ProducerMessage, bufferSize)
	sp.stop = make(chan struct{})
	sp.producer, err = sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}
	go sp.daemonProducer()
	log.Printf("[syncProducer] start, brokers:%v, bufferSize:%v", brokers, bufferSize)
	return sp, nil
}

// Send use sync producer
func (sp *SyncProducer) Send(msg *sarama.ProducerMessage) {
	if msg != nil {
		sp.messages <- msg
	}
}

// daemon send msg to special topic with sync producer
func (sp *SyncProducer) daemonProducer() {
	for {
		mes, ok := <-sp.messages
		if !ok {
			sp.stop <- struct{}{}
			return
		}
		m, _ := utils.ToJsonString(mes)
		partition, offset, err := sp.producer.SendMessage(mes)
		if err != nil {
			log.Printf("[syncProducer] return, partition:%v, offset:%v, msg:%v, err:%v",
				partition, offset, m, err.Error())
		}
	}
}

// Close sync producer
func (sp *SyncProducer) Close() error {
	close(sp.messages)
	<-sp.stop
	return sp.producer.Close()
}
