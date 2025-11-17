package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/jimtang2/bc-go/lib/config"
	"github.com/jimtang2/bc-go/lib/db"
	"github.com/jimtang2/bc-go/lib/kafka"
	"github.com/jimtang2/bc-go/lib/stream"
	"github.com/jimtang2/bc-go/lib/stream/driver"
	"github.com/spf13/viper"
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

func init() {
	viper.SetDefault("http.port", ":8080") // for healthcheck
}

func main() {
	config.Read()
	tickersProcessor := NewTickersProcessor()
	pairsByExchange, err := db.GetTrackedPairs()
	if err != nil {
		log.Fatal(err)
	}
	producer, err := kafka.NewAsyncProducer()
	if err != nil {
		log.Fatal(err)
	}
	//
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
	// healthcheck notes:
	// if tickers processing fails, it means no match is found for some amount of time, which is not expected within any sufficient period of time (> ~1-3s), in turn it means all of the stream drivers may be failing to stream tickers data; so the logic is each individual stream driver should auto restart on their own before tickers processing is deemed to have failed after x amount of time in order to deem service to be in failed state
	go func() {
		log.Printf("[%v -> %v/%v] on", topics[0], topics[1], topics[2])
		kafka.ProcessStream(kafka.ProcessorConfig{
			Topic:     topics[0],
			Offset:    sarama.OffsetNewest,
			Processor: tickersProcessor, // computes spreads and alpha from Topic/Offset
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

	http.Handle("/health", HealthHandler(tickersProcessor))
	if err := http.ListenAndServe(viper.GetString("http.port"), nil); err != nil {
		log.Fatal(err)
	}
}
