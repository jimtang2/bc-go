package stream

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/jimtang2/bc-go/lib/db"
)

type Coinbase struct{}

func (s *Coinbase) Start() chan Message {
	c := make(chan Message)
	pairs := db.TrackedPairs("coinbase")
	url := "wss://ws-feed.exchange.coinbase.com"
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("[coinbase]", err)
		return c
	}
	subscribeMsg := map[string]interface{}{
		"type":        "subscribe",
		"product_ids": pairs,
		"channels":    []string{"ticker"},
	}
	if err := ws.WriteJSON(subscribeMsg); err != nil {
		log.Println(err)
		close(c)
		return c
	}
	log.Println("[coinbase] subscribed to", strings.Join(pairs, ", "))
	go func() {
		defer ws.Close()
		for {
			messageType, b, err := ws.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			switch messageType {
			case websocket.PingMessage:
			case websocket.TextMessage:
				m, err := parseCoinbaseTickerMessage(b)
				if err != nil {
					log.Println(err)
					continue
				}
				c <- Message{
					Topic: "tickers",
					Key:   "[coinbase]" + m.ProductID,
					Headers: map[string]string{
						"exchange": "coinbase",
						"api":      "ticker channel",
						"pair":     m.ProductID,
					},
					Payload: b,
				}
			default:
			}
		}
	}()
	return c
}

// https://docs.cdp.coinbase.com/exchange/websocket-feed/channels#ticker-channel
type CoinbaseTickerMessageMinimal struct {
	ProductID string `json:"product_id"`
}

func parseCoinbaseTickerMessage(b []byte) (CoinbaseTickerMessageMinimal, error) {
	v := CoinbaseTickerMessageMinimal{}
	err := json.Unmarshal(b, &v)
	return v, err
}
