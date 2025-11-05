/*
The main function of this package is ProcessStream(ProcessorConfig) which provides a high level API to process records from any topic
*/
package kafka

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("kafka.brokers", []string{"0.0.0.0:29092"})
}

func NewAsyncProducer() (sarama.AsyncProducer, error) {
	var (
		brokers = viper.GetStringSlice("kafka.brokers")
		config  = sarama.NewConfig()
	)
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = false
	return sarama.NewAsyncProducer(brokers, config)
}

type Processor interface {
	Process(m *sarama.ConsumerMessage) KMessage
}

type ProcessorConfig struct {
	Topic      string
	Offset     int64
	Processor  Processor
	OnKMessage func(kmessage KMessage)
}

type KMessage interface {
	Key() string
	Bytes() []byte
	IsValid() bool
}

func ProcessStream(processorConfig ProcessorConfig) {
	var (
		brokers         = viper.GetStringSlice("kafka.brokers")
		config          = sarama.NewConfig()
		topic           = processorConfig.Topic
		offset          = processorConfig.Offset
		processor       = processorConfig.Processor
		kmessageHandler = processorConfig.OnKMessage
	)
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer consumer.Close()
		partitionList, err := consumer.Partitions(topic)
		if err != nil {
			log.Fatal(err)
		}
		for _, partition := range partitionList {
			pc, err := consumer.ConsumePartition(topic, partition, offset)
			if err != nil {
				log.Fatal(err)
			}
			go func(pc sarama.PartitionConsumer) {
				defer pc.Close()
				for {
					kmessageHandler(processor.Process(<-pc.Messages()))
				}
			}(pc)
		}
		select {}
	}()
}
