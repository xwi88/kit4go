// Package kafka producer, sync producer
package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/kdpujie/log4go"

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
	addrs    []string
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
	sp.cfg = cfg
	sp.messages = make(chan *sarama.ProducerMessage, bufferSize)
	sp.stop = make(chan struct{})
	sp.addrs = brokers
	sp.producer, err = sarama.NewSyncProducer(brokers, sp.cfg)
	if err != nil {
		return nil, err
	}
	go sp.daemonProducer()
	log4go.Debug("[syncProducer] created, brokers:%v, bufferSize:%v", brokers, bufferSize)
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
			log4go.Error("[syncProducer] return, partition:%v, offset:%v, msg:%v, err:%v",
				partition, offset, m, err.Error())
		}
	}
}

// Close sync producer
func (sp *SyncProducer) Close() error {
	close(sp.messages)
	<-sp.stop
	log4go.Info("[syncProducer] close, brokers:%v, bufferSize:%v", sp.addrs, sp.cfg.ChannelBufferSize)
	return sp.producer.Close()
}
