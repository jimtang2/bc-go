// webhook$$bc-go;cmd/trade/main.go;grok$$
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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
		log.Println("config not found (using defaults)")
		log.Println(viper.AllSettings())
	}
}

func main() {
	log.Println("trade on")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
