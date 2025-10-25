// webhook$$bc-go;cmd/stream/proxy.go;grok$$
package main

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/jimtang2/bc-go/lib/stream"
	"github.com/spf13/viper"
)

type Proxy struct {
	producer sarama.SyncProducer
	channel  chan stream.Message
	streams  []stream.Stream
}

func NewProxy() (p *Proxy, err error) {
	p = &Proxy{
		channel: make(chan stream.Message),
		streams: []stream.Stream{
			&stream.Binance{},
			&stream.Bitfinex{},
			&stream.Coinbase{},
			&stream.Kraken{},
			&stream.OKX{},
			// &stream.Gemini{},
			// &stream.Uniswap{},
		},
	}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	if p.producer, err = sarama.NewSyncProducer(
		viper.GetStringSlice("kafka.brokers"),
		config,
	); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Proxy) write(m stream.Message) error {
	message := &sarama.ProducerMessage{
		Topic: m.Topic,
		Key:   sarama.StringEncoder(m.Key),
		Value: sarama.ByteEncoder(m.Payload),
	}
	if m.Headers != nil {
		for k, v := range m.Headers {
			message.Headers = append(message.Headers, sarama.RecordHeader{
				Key:   sarama.ByteEncoder(k),
				Value: sarama.ByteEncoder(v),
			})
		}
	}
	_, _, err := p.producer.SendMessage(message)
	return err
}

func (p *Proxy) Start() {
	go func() {
		for {
			if err := p.write(<-p.channel); err != nil {
				log.Fatal(err)
			}
		}
	}()
	for _, s := range p.streams {
		go func(s stream.Stream) {
			ch := s.Start()
			for {
				p.channel <- <-ch
			}
		}(s)
	}
}
