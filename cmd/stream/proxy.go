// webhook$$bc-go;cmd/stream/proxy.go;grok$$
package main

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/jimtang2/bc-go/lib/stream"
	"github.com/spf13/viper"
)

// proxy calls Start on a set of stream.Stream interfaces and collects stream.Message's on its channel to publish to kafka topics.
type Proxy struct {
	producer sarama.SyncProducer
	channel  chan stream.Message
	streams  []stream.Stream
}

func NewProxy() (*Proxy, error) {
	var (
		err     error
		brokers = viper.GetStringSlice("kafka.brokers")
		config  = sarama.NewConfig()
		p       = Proxy{
			channel: make(chan stream.Message),
			streams: []stream.Stream{
				&stream.BinanceStream{},
				&stream.OKXStream{},
				&stream.BitfinexStream{},
				&stream.CoinbaseStream{},
				&stream.KrakenStream{},
				// &stream.Gemini{},
				// &stream.Uniswap{},
			},
		}
	)
	config.Producer.Return.Successes = true
	p.producer, err = sarama.NewSyncProducer(brokers, config)
	return &p, err
}

func (p *Proxy) Produce() {
	for {
		_, _, err := p.producer.SendMessage(p.NewMessage(<-p.channel))
		if err != nil {
			log.Println("[proxy]", err)
		}
	}
}

func (p *Proxy) Stream() {
	for _, s := range p.streams {
		go func(s stream.Stream) {
			ch := s.Start()
			for {
				p.channel <- <-ch
			}
		}(s)
	}
	select {}
}

func (p *Proxy) NewMessage(sm stream.Message) *sarama.ProducerMessage {
	m := &sarama.ProducerMessage{
		Topic: sm.Topic,
		Key:   sarama.StringEncoder(sm.Key),
		Value: sarama.ByteEncoder(sm.Payload),
	}
	if sm.Headers != nil {
		for k, v := range sm.Headers {
			m.Headers = append(m.Headers, sarama.RecordHeader{
				Key:   sarama.ByteEncoder(k),
				Value: sarama.ByteEncoder(v),
			})
		}
	}
	return m
}
