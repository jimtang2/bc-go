package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/jimtang2/bc-go/lib/config"
	"github.com/jimtang2/bc-go/lib/kafka"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("http.port", ":8080")
	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:5173"})
}

func main() {
	config.Read()
	socketHandler := NewSocketHandler()
	kafka.ProcessStream(kafka.ProcessorConfig{
		Topic:      "spreads",
		Offset:     sarama.OffsetNewest,
		Processor:  NewSocketProxyProcessor(socketHandler.channel),
		OnKMessage: func(kmessage kafka.KMessage) {},
	})
	http.Handle("/", http.FileServer(http.Dir("dist")))
	http.Handle("/alpha/", http.StripPrefix("/alpha", http.FileServer(http.Dir("dist"))))
	http.Handle("/ws", socketHandler)
	http.Handle("/api/alpha", &AlphaHandler{mu: sync.Mutex{}})
	log.Println("dashboard listening on", viper.GetString("http.port"))
	if err := http.ListenAndServe(viper.GetString("http.port"), cors.New(cors.Options{
		AllowedOrigins:   viper.GetStringSlice("cors.allowed_origins"),
		AllowedMethods:   []string{"GET", "POST"},
		AllowCredentials: false,
	}).Handler(http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
