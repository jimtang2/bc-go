// webhook$$bc-go;dashboard/server/main.go;grok$$
package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")
	viper.SetDefault("http.port", ":8080")
	viper.SetDefault("kafka.brokers", []string{"0.0.0.0:29092"})
	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:5173"})
	if err := viper.ReadInConfig(); err != nil {
		log.Println("config not found (using defaults)")
		log.Println(viper.AllSettings())
	}
}

var channel = make(chan []byte)

func main() {
	go startHttp()
	go consume()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

func startHttp() {
	http.Handle("/", http.FileServer(http.Dir("dist")))
	http.Handle("/ws", NewSocketHandler())
	log.Println("dashboard listening on", viper.GetString("http.port"))
	if err := http.ListenAndServe(viper.GetString("http.port"), cors.New(cors.Options{
		AllowedOrigins:   viper.GetStringSlice("cors.allowed_origins"),
		AllowedMethods:   []string{"GET", "POST"},
		AllowCredentials: false,
	}).Handler(http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
}

func consume() {
	var (
		err           error
		brokers       = viper.GetStringSlice("kafka.brokers")
		config        = sarama.NewConfig()
		topicIn       = "alpha"
		initialOffset = -200
		consumer      sarama.Consumer
		processor     = NewProcessor()
	)
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	if consumer, err = sarama.NewConsumer(brokers, config); err != nil {
		log.Fatal(err)
	}
	go func() {
		defer consumer.Close()
		partitionList, err := consumer.Partitions(topicIn)
		if err != nil {
			log.Fatal(err)
		}
		for _, partition := range partitionList {
			latest, err := consumer.GetOffset(topicIn, partition, sarama.OffsetNewest)
			if err != nil {
				log.Fatal(err)
			}
			target := latest + initialOffset
			if target < 0 {
				target = sarama.OffsetOldest
			}
			pc, err := consumer.ConsumePartition(topicIn, partition, target)
			if err != nil {
				log.Fatal(err)
			}
			go func(pc sarama.PartitionConsumer) {
				defer pc.Close()
				for {
					if m, ok := <-pc.Messages(); ok {
						if b, err := processor.process(m); err == nil {
							channel <- b
						}
					}
				}
			}(pc)
		}
		select {}
	}()
	select {}
}
