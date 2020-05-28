// Package kafka producer, async producer
package kafka

import (
	"log"

	"github.com/Shopify/sarama"
)

var (
	// DefaultAsyncProducerBufferSize default async buffer size
	DefaultAsyncProducerBufferSize = 1
)

// AsyncProducer async producer
type AsyncProducer struct {
	producer sarama.AsyncProducer
	messages chan *sarama.ProducerMessage
	cfg      *sarama.Config
	stop     chan struct{}
}

// NewAsyncProducer create async producer instance
func NewAsyncProducer(brokers []string, bufferSize int, cfg *sarama.Config) (ap *AsyncProducer, err error) {
	ap = new(AsyncProducer)
	// if not set, use kafka default ChannelBufferSize=256
	if bufferSize == 0 {
		bufferSize = cfg.ChannelBufferSize
	} else if bufferSize <= DefaultAsyncProducerBufferSize {
		// if set, but not greater than our limit DefaultAsyncProducerBufferSize, set it
		bufferSize = DefaultAsyncProducerBufferSize
	}
	ap.cfg = cfg
	ap.messages = make(chan *sarama.ProducerMessage, bufferSize)
	ap.stop = make(chan struct{})
	ap.producer, err = sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}
	go ap.daemonProducer()
	log.Printf("[asyncProducer] start, brokers:%v, bufferSize:%v", brokers, bufferSize)
	return ap, nil
}

// Send use async producer
func (ap *AsyncProducer) Send(msg *sarama.ProducerMessage) {
	if msg != nil {
		ap.messages <- msg
	}
}

// daemon send msg to special topic with async producer
func (ap *AsyncProducer) daemonProducer() {
	// consume successes
	go func() {
		if ap.cfg.Producer.Return.Successes {
			for pm := range ap.producer.Successes() {
				log.Printf("[asyncProducer] return, successes:%v", pm)
			}
		}
	}()

	// consume errors
	go func() {
		if ap.cfg.Producer.Return.Errors {
			for pe := range ap.producer.Errors() {
				log.Printf("[asyncProducer] return, errors:%v", pe.Error())
			}
		}
	}()

	for {
		mes, ok := <-ap.messages
		if !ok {
			ap.stop <- struct{}{}
			return
		}
		ap.producer.Input() <- mes
	}
}

// Close async producer
func (ap *AsyncProducer) Close() error {
	close(ap.messages)
	<-ap.stop
	return ap.producer.Close()
}
