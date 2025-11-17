package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/jimtang2/bc-go/lib/config"
	"github.com/jimtang2/bc-go/lib/kafka"
	"github.com/rs/cors"
	"github.com/spf13/viper"

	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

//go:embed dist/**
var staticFiles embed.FS

func init() {
	viper.SetDefault("http.port", ":8080")
	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:5173"})
	viper.SetDefault("healthcheck.last_rcv_s", 10)
}

func main() {
	config.Read()
	fs := statigz.FileServer(staticFiles, brotli.AddEncoding, statigz.FSPrefix("dist"))
	matchesHandler := NewMatchesHandler()
	kafka.ProcessStream(kafka.ProcessorConfig{
		Topic:      "stream_matches",
		Offset:     sarama.OffsetNewest,
		Processor:  NewMatchesProxyProcessor(matchesHandler.channel),
		OnKMessage: func(kmessage kafka.KMessage) {},
	})
	tickersHandler := NewTickersHandler()
	kafka.ProcessStream(kafka.ProcessorConfig{
		Topic:      "stream_tickers",
		Offset:     sarama.OffsetNewest,
		Processor:  NewTickersProxyProcessor(tickersHandler.channel),
		OnKMessage: func(kmessage kafka.KMessage) {},
	})
	alphaHandler := &AlphaHandler{}
	exchangesHandler := &ExchangesHandler{}
	http.Handle("/", fs)
	http.Handle("/alpha/", http.StripPrefix("/alpha", fs))
	http.Handle("/ws/matches", matchesHandler)
	http.Handle("/ws/tickers", tickersHandler)
	http.Handle("/api/alpha", alphaHandler)
	http.Handle("/api/exchanges", exchangesHandler)
	http.Handle("/health", HealthHandler(matchesHandler, tickersHandler, alphaHandler, exchangesHandler))
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
