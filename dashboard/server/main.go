// webhook$$bc-go;dashboard/server/main.go;grok$$
package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/cors"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")
	viper.SetDefault("http.port", ":8090")
	viper.SetDefault("kafka.brokers", []string{"0.0.0.0:29092"})
	viper.SetDefault("db.url", "postgres://postgres:password@localhost/bc?sslmode=disable")
	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:5173"})

	if err := viper.ReadInConfig(); err != nil {
		log.Println("config not found (using defaults)")
		log.Println(viper.AllSettings())
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("dist")))
	http.Handle("/ws", NewSocketHandler())
	http.HandleFunc("/lists", http.HandlerFunc(serveListsFunc))
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
