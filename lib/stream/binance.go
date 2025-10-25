package stream

import (
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/jimtang2/bc-go/lib/db"
	"github.com/jimtang2/bc-go/lib/messaging"
)

type Binance struct{}

func (s *Binance) Start() chan Message {
	c := make(chan Message)
	pairs := db.TrackedPairs("binance")
	url := fmt.Sprintf("wss://stream.binance.com:9443/stream?streams=%s@ticker", strings.Join(pairs, "@ticker/"))
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("[binance]", err)
		return c
	}
	log.Println("[binance] subscribed to", strings.Join(pairs, ", "))
	go func() {
		defer ws.Close()
		for {
			messageType, b, err := ws.ReadMessage()
			if err != nil {
				log.Println("[binance]", err)
				return
			}
			switch messageType {
			case websocket.PingMessage:
				err = ws.WriteMessage(websocket.PongMessage, b)
				if err != nil {
					log.Println("[binance]", err)
					return
				}
				continue
			case websocket.TextMessage:
				key, err := messaging.BinanceTickerCombinedStream(b)
				if err != nil {
					log.Println("[binance]", err)
					continue
				}
				c <- Message{
					Topic: "tickers",
					Key:   "[binance]" + key,
					Headers: map[string]string{
						"url":      url,
						"exchange": "binance",
						"api":      "ticker combined stream",
					},
					Payload: b,
				}
			default:
			}
		}
	}()
	return c
}
