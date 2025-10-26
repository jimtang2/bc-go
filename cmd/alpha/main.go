// webhook$$bc-go;cmd/alpha/main.go;grok$$
package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/jimtang2/bc-go/lib/stream"
	"github.com/spf13/viper"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")
	viper.SetDefault("kafka.brokers", []string{"0.0.0.0:29092"})
	if err := viper.ReadInConfig(); err != nil {
		log.Println("applying default config")
		log.Println(viper.AllSettings())
	}
}

func main() {
	var (
		err      error
		consumer sarama.Consumer
		producer sarama.SyncProducer
		parser   = &stream.Parser{}
		channel  = make(chan *sarama.ConsumerMessage)
		topicIn  = "tickers_raw"
		topicOut = "tickers"
		brokers  = viper.GetStringSlice("kafka.brokers")
		config   = sarama.NewConfig()
	)
	config.Producer.Return.Successes = true
	if consumer, err = sarama.NewConsumer(brokers, config); err != nil {
		log.Fatal(err)
	}
	if producer, err = sarama.NewSyncProducer(brokers, config); err != nil {
		log.Fatal(err)
	}
	go func() {
		defer consumer.Close()
		partitionList, err := consumer.Partitions(topicIn)
		if err != nil {
			log.Fatal(err)
			return
		}
		for _, partition := range partitionList {
			pc, err := consumer.ConsumePartition(topicIn, partition, sarama.OffsetNewest)
			if err != nil {
				log.Fatal(err)
			}
			go func(pc sarama.PartitionConsumer) {
				for {
					channel <- <-pc.Messages()
				}
			}(pc)
		}
		for {
			t := parser.Parse(<-channel)
			if t == nil {
				continue
			}
			b, _ := json.Marshal(t)
			message := &sarama.ProducerMessage{
				Topic: topicOut,
				Key:   sarama.StringEncoder(t.Key()),
				Value: sarama.ByteEncoder(b),
			}
			if _, _, err := producer.SendMessage(message); err != nil {
				log.Println(err)
			}
		}
	}()
	log.Println("alpha on")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
