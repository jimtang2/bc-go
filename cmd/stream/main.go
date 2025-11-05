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
	_ "github.com/jimtang2/bc-go/lib/stream/driver"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	config.Read()
	pairsByExchange, err := db.GetTrackedPairs()
	if err != nil {
		log.Fatal(err)
	}
	producer, err := kafka.NewAsyncProducer()
	if err != nil {
		log.Fatal(err)
	}
	drivers := []string{
		"binance-tickers",
		"bitfinex-tickers",
		"coinbase-tickers",
		"kraken-tickers",
		"okx-tickers",
	}
	go func() {
		for _, driverName := range drivers {
			driverConfig := stream.DriverConfig{
				Args: pairsByExchange[strings.Split(driverName, "-")[0]],
				OnKMessage: func(kmessage kafka.KMessage) {
					producer.Input() <- &sarama.ProducerMessage{
						Topic: "tickers",
						Key:   sarama.StringEncoder(kmessage.Key()),
						Value: sarama.ByteEncoder(kmessage.Bytes()),
					}
				},
				OnError: func(d stream.Driver, cfg stream.DriverConfig, err error) {
					log.Println("restarting", driverName)
					if err := d.Close(); err != nil {
						log.Println(err)
					}
					time.Sleep(3 * time.Second)
					stream.Open(driverName, cfg)
				},
			}
			stream.Open(driverName, driverConfig)
		}
	}()
	log.Println("streaming")
	// spreads and edges
	kafka.ProcessStream(kafka.ProcessorConfig{
		Topic:     "tickers",
		Offset:    sarama.OffsetNewest,
		Processor: NewTickersProcessor(), // computes spreads and alpha
		OnKMessage: func(kmessage kafka.KMessage) {
			a, ok := kmessage.(*Alpha)
			if !ok || a == nil || !a.IsValid() {
				return
			}
			if a.Profit > 0 {
				producer.Input() <- &sarama.ProducerMessage{
					Topic: "alpha",
					Key:   sarama.StringEncoder(a.Key()),
					Value: sarama.ByteEncoder(a.Bytes()),
				}
			}
			producer.Input() <- &sarama.ProducerMessage{
				Topic: "spreads",
				Key:   sarama.StringEncoder(a.Key()),
				Value: sarama.ByteEncoder(a.Bytes()),
			}
		},
	})
	log.Println("processing tickers")
	// spreads and edges
	kafka.ProcessStream(kafka.ProcessorConfig{
		Topic:      "alpha",
		Offset:     sarama.OffsetOldest,
		Processor:  NewAlphaProcessor(), // inserts new alpha to db
		OnKMessage: func(kmessage kafka.KMessage) {},
	})
	log.Println("processing alpha")
	select {}
}
