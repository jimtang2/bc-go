package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
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

func main() {
	socketHandler := NewSocketHandler()
	go consume("spreads", sarama.OffsetNewest, NewSpreadsProcessor(socketHandler.channel))
	go consume("alpha", sarama.OffsetOldest, NewAlphaProcessor())
	go func() {
		http.Handle("/", http.FileServer(http.Dir("dist")))
		http.Handle("/alpha/", http.StripPrefix("/alpha", http.FileServer(http.Dir("dist"))))
		http.Handle("/ws", socketHandler)
		http.Handle("/api/alpha", &AlphaHandler{mu: sync.Mutex{}})
		log.Println("dashboard listening on", viper.GetString("http.port"))
		handler := cors.New(cors.Options{
			AllowedOrigins:   viper.GetStringSlice("cors.allowed_origins"),
			AllowedMethods:   []string{"GET", "POST"},
			AllowCredentials: false,
		}).Handler(http.DefaultServeMux)
		if err := http.ListenAndServe(viper.GetString("http.port"), handler); err != nil {
			log.Fatal(err)
		}
	}()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

func consume(topic string, offset int64, processor Processor) {
	var (
		err      error
		brokers  = viper.GetStringSlice("kafka.brokers")
		config   = sarama.NewConfig()
		consumer sarama.Consumer
	)
	config.Consumer.Offsets.Initial = offset
	if consumer, err = sarama.NewConsumer(brokers, config); err != nil {
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
					if m, ok := <-pc.Messages(); ok {
						if err := processor.Process(m); err != nil {
							log.Println(err)
							continue
						}
					}
				}
			}(pc)
		}
		select {}
	}()
	select {}
}
