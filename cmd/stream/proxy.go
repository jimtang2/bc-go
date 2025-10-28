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
				stream.NewBinanceStream(),
				stream.NewOKXStream(),
				stream.NewBitfinexStream(),
				stream.NewCoinbaseStream(),
				stream.NewKrakenStream(),
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
		var (
			sm = <-p.channel // receive stream.Message
			m  = &sarama.ProducerMessage{
				Topic: sm.Topic,
				Key:   sarama.StringEncoder(sm.Key),
				Value: sarama.ByteEncoder(sm.Payload),
			}
		)
		if sm.Headers != nil {
			for k, v := range sm.Headers {
				m.Headers = append(m.Headers, sarama.RecordHeader{
					Key:   sarama.ByteEncoder(k),
					Value: sarama.ByteEncoder(v),
				})
			}
		}
		_, _, err := p.producer.SendMessage(m)
		if err != nil {
			log.Println("[proxy]", err)
		}
	}
}

func (p *Proxy) Stream() {
	for _, s := range p.streams {
		go func(s stream.Stream) {
			go s.Start()
			for {
				p.channel <- <-s.Output()
			}
		}(s)
	}
	select {}
}
