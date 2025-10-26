// webhook$$bc-go;cmd/stream/main.go;grok$$
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")
	viper.SetDefault("kafka.brokers", []string{"0.0.0.0:29092"})
	if err := viper.ReadInConfig(); err != nil {
		log.Println("config not found (using defaults)")
		log.Println(viper.AllSettings())
	}
}

func main() {
	proxy, err := NewProxy()
	if err != nil {
		log.Fatal(err)
	}
	go proxy.Stream()
	go proxy.Produce()
	log.Println("stream on")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
