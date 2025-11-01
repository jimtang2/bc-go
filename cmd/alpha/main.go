// webhook$$bc-go;cmd/alpha/main.go;grok$$
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func init() {
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
		err       error
		brokers   = viper.GetStringSlice("kafka.brokers") // kafka brokers
		config    = sarama.NewConfig()                    // kafka config
		topicIn   = "tickers"                             // consumer topic
		topicsOut = []string{"spreads", "alpha"}          // producer topic
		consumer  sarama.Consumer                         // kafka consumer
		producer  sarama.SyncProducer                     // kafka producer
		processor = NewProcessor()                        // stream processor
	)
	config.Producer.Return.Successes = true
	if consumer, err = sarama.NewConsumer(brokers, config); err != nil {
		log.Fatal(err)
	}
	if producer, err = sarama.NewSyncProducer(brokers, config); err != nil {
		log.Fatal(err)
	}
	// consumption
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
				defer pc.Close()
				for {
					processor.Process(<-pc.Messages())
				}
			}(pc)
		}
		select {}
	}()
	// production
	go func() {
		for {
			a := <-processor.alphaChan
			go func(a *Alpha) {
				if _, _, err := producer.SendMessage(&sarama.ProducerMessage{
					Topic: topicsOut[0],
					Key:   sarama.StringEncoder(a.Key()),
					Value: sarama.ByteEncoder(a.Bytes()),
				}); err != nil {
					log.Println(err)
				}
			}(a)
			go func(a *Alpha) {
				if a.Profit <= 0 {
					return
				}
				if _, _, err := producer.SendMessage(&sarama.ProducerMessage{
					Topic: topicsOut[1],
					Key:   sarama.StringEncoder(a.Key()),
					Value: sarama.ByteEncoder(a.Bytes()),
				}); err != nil {
					log.Println(err)
				}
			}(a)
		}
	}()
	log.Println("alpha on")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
