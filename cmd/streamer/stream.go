// webhook$$bc-go;cmd/streamer/stream.go;grok$$
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var db *sql.DB

func trackedPairs() (l []string) {
	var (
		err error
		b   []byte
	)
	if db == nil {
		db, err = sql.Open("postgres", viper.GetString("db.url"))
		if err != nil {
			log.Println(err)
			return l
		}
	}
	if err := db.QueryRow(`select array_to_json(array_agg(name)) from pairs`).Scan(&b); err != nil {
		log.Println(err)
		return l
	}
	if err := json.Unmarshal(b, &l); err != nil {
		log.Println(err)
		return l
	}
	return l
}

type Stream interface {
	Start() chan Message
}

type Binance struct{}

func (s *Binance) Start() chan Message {
	var (
		c     = make(chan Message)
		pairs = trackedPairs()
		url   = fmt.Sprintf("wss://stream.binance.com:9443/stream?streams=%s@ticker", strings.Join(pairs, "@ticker/"))
	)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer conn.Close()
		for {
			messageType, b, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			if messageType == websocket.PingMessage {
				err = conn.WriteMessage(websocket.PongMessage, b)
				if err != nil {
					log.Println(err)
					return
				}
				continue
			}
			c <- Message{
				Key: url,
				Headers: map[string]string{
					"exchange": "binance",
					"api":      "combined ticker stream",
				},
				Payload: b,
			}
		}
	}()
	return c
}

type Bitfinex struct{}

func (s *Bitfinex) Start() chan Message {
	c := make(chan Message)

	return c
}

type Uniswap struct{}

func (s *Uniswap) Start() chan Message {
	c := make(chan Message)

	return c
}

type Coinbase struct{}

func (s *Coinbase) Start() chan Message {
	c := make(chan Message)

	return c
}

type Kraken struct{}

func (s *Kraken) Start() chan Message {
	c := make(chan Message)

	return c
}
