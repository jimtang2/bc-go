package main

import (
	"log"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/jimtang2/bc-go/lib/config"
	"github.com/jimtang2/bc-go/lib/db"
	"github.com/jimtang2/bc-go/lib/kafka"
	"github.com/jimtang2/bc-go/lib/stream"
	"github.com/jimtang2/bc-go/lib/stream/driver"
)

var (
	restartInterval = 3 * time.Second
	drivers         = []string{
		"binance-tickers",
		"bitfinex-tickers",
		"coinbase-tickers",
		"kraken-tickers",
		"okx-tickers",
	}
	topics = []string{
		"stream_tickers",
		"stream_alpha",
		"stream_matches",
	}
)

func main() {
	config.Read()
	pairsByExchange, err := db.GetTrackedPairs()
	if err != nil {
		log.Fatal(err)
	}
	producer, err := kafka.NewAsyncProducer()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		log.Printf("[stream -> %v] on", topics[0])
		for _, driverName := range drivers {
			driverConfig := stream.DriverConfig{
				Args: pairsByExchange[strings.Split(driverName, "-")[0]],
				OnKMessage: func(kmessage kafka.KMessage) {
					m, ok := kmessage.(*driver.Ticker)
					if ok && m.IsValid() {
						producer.Input() <- &sarama.ProducerMessage{
							Topic: topics[0],
							Key:   sarama.StringEncoder(kmessage.Key()),
							Value: sarama.ByteEncoder(kmessage.Bytes()),
						}
					}
				},
				OnError: func(d stream.Driver, cfg stream.DriverConfig, err error) {
					log.Println("restarting", driverName)
					if err := d.Close(); err != nil {
						log.Println(err)
					}
					time.Sleep(restartInterval)
					stream.Open(driverName, cfg)
				},
			}
			stream.Open(driverName, driverConfig)
		}
	}()
	go func() {
		log.Printf("[%v -> %v/%v] on", topics[0], topics[1], topics[2])
		kafka.ProcessStream(kafka.ProcessorConfig{
			Topic:     topics[0],
			Offset:    sarama.OffsetNewest,
			Processor: NewTickersProcessor(), // computes spreads and alpha
			OnKMessage: func(kmessage kafka.KMessage) {
				m, ok := kmessage.(*Match)
				if !ok || m == nil {
					return
				}
				if m.Match.Calculations.ProfitLoss > 0 {
					producer.Input() <- &sarama.ProducerMessage{
						Topic: topics[1],
						Key:   sarama.StringEncoder(m.Key()),
						Value: sarama.ByteEncoder(m.Bytes()),
					}
				}
				producer.Input() <- &sarama.ProducerMessage{
					Topic: topics[2],
					Key:   sarama.StringEncoder(m.Key()),
					Value: sarama.ByteEncoder(m.Bytes()),
				}
			},
		})
	}()
	go func() {
		log.Printf("[%v -> sql] on", topics[1])
		kafka.ProcessStream(kafka.ProcessorConfig{
			Topic:      topics[1],
			Offset:     sarama.OffsetOldest,
			Processor:  &AlphaProcessor{},
			OnKMessage: func(kmessage kafka.KMessage) {},
		})
	}()
	select {}
}
