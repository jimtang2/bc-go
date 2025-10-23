// webhook$$bc-go;cmd/analyzer/main.go;grok$$
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
)

// init function for global scope configu
func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// orchestration of main app objects
// set up kafka client aka Producer
// the producer simply listens to quote messages on a channel that is passed to all the quote fetchers and pushes the messages on the quotes topic
func main() {
	c := NewConsumer()
	go c.Consume()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
