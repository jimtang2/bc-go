// webhook$$bc-go;cmd/print_proto/main.go;grok$$
package main

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/jimtang2/bc-go/lib/config"
	_ "github.com/jimtang2/bc-go/lib/kafka"
	"github.com/jimtang2/bc-go/pkg/pb/v1"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"
)

func main() {
	config.Read()

	offsetIDs := []int64{
		5253369,
		5253363,
		5253357,
	}
	topic := "stream_matches"
	brokers := viper.GetStringSlice("kafka.brokers")

	// Create client to query metadata
	client, err := sarama.NewClient(brokers, nil)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}
	defer client.Close()

	// Create consumer for actual message fetching
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Fatal("Failed to create consumer:", err)
	}
	defer consumer.Close()

	// Get all partitions
	partitions, err := client.Partitions(topic)
	if err != nil {
		log.Fatal("Failed to get partitions:", err)
	}

	for _, targetOffset := range offsetIDs {
		var foundPartition int32
		var found bool

		// Find which partition contains this offset
		for _, partition := range partitions {
			oldest, err := client.GetOffset(topic, partition, sarama.OffsetOldest)
			if err != nil {
				continue
			}
			newest, err := client.GetOffset(topic, partition, sarama.OffsetNewest)
			if err != nil {
				continue
			}

			if targetOffset >= oldest && targetOffset <= newest {
				foundPartition = partition
				found = true
				break
			}
		}

		if !found {
			log.Printf("Offset %d not found in any partition", targetOffset)
			continue
		}

		// Consume from exact offset
		pc, err := consumer.ConsumePartition(topic, foundPartition, targetOffset)
		if err != nil {
			log.Printf("Failed to start consumer for offset %d: %v", targetOffset, err)
			continue
		}

		// Read one message
		select {
		case msg := <-pc.Messages():

			match := pb.Match{}
			proto.Unmarshal(msg.Value, &match)
			log.Println(match)

			// log.Printf("SUCCESS: offset=%d partition=%d key=%s value=%s", msg.Offset, msg.Partition, string(msg.Key), string(msg.Value))
		case err := <-pc.Errors():
			log.Printf("Error reading offset %d: %v", targetOffset, err)
		}

		pc.AsyncClose()
	}
}
