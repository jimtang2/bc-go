// webhook$$bc-go;dashboard/server/dashboard.go;grok$$
package main

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

type Consumer struct{}

func (s *Consumer) consume(in chan interface{}) {
	var (
		brokers = viper.GetStringSlice("kafka.brokers")
		config  = sarama.NewConfig()
	)
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Println(err)
		close(in)
		return
	}
	defer consumer.Close()
	topic := "tickers"
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		log.Println(err)
		close(in)
		return
	}
	for _, partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Println(err)
			continue
		}
		go func(pc sarama.PartitionConsumer) {
			defer pc.Close()
			for {
				in <- string((<-pc.Messages()).Value)
			}
		}(pc)
	}
	select {}
}
