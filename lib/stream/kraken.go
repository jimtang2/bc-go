package stream

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jimtang2/bc-go/lib/db"
)

type Kraken struct {
	pairs map[int]KrakenSubscriptionMessage
	mu    sync.Mutex
}

// https://docs.kraken.com/api/docs/websocket-v2/ticker
func (s *Kraken) Start() chan Message {
	s.pairs = map[int]KrakenSubscriptionMessage{}
	s.mu = sync.Mutex{}
	c := make(chan Message)
	pairs := db.TrackedPairs("kraken")
	url := "wss://ws.kraken.com"
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("[kraken]", err)
		return c
	}
	subscribeMsg := map[string]interface{}{
		"event": "subscribe",
		"pair":  pairs, // e.g., ["XBT/USD", "ETH/USD"]
		"subscription": map[string]string{
			"name": "ticker",
		},
	}
	ws.WriteJSON(subscribeMsg)
	go func() {
		defer ws.Close()
		for {
			messageType, b, err := ws.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			switch messageType {
			case websocket.TextMessage:
				firstChar := string(b[0])
				switch firstChar {
				case "{":
					m, err := parseKrakenSubscriptionMessage(b)
					if err != nil {
						log.Println(err)
						continue
					}
					switch m.Event {
					case "heartbeat":
					case "subscriptionStatus":
						if m.Status == "subscribed" {
							s.mu.Lock()
							s.pairs[m.ChannelID] = m
							s.mu.Unlock()
						}
						log.Println("[kraken] subscribed to", m.Pair)
					default:
					}
				case "[":
					chanID, _, err := parseKrakenTickerMessage(b)
					if err != nil {
						log.Println(err)
						continue
					}
					s.mu.Lock()
					m := s.pairs[chanID]
					s.mu.Unlock()
					c <- Message{
						Topic: "tickers",
						Key:   "[kraken]" + m.Pair,
						Headers: map[string]string{
							"exchange":    "kraken",
							"api":         "ticker subscription",
							"channelID":   fmt.Sprintf("%v", m.ChannelID),
							"channelName": m.ChannelName,
							"pair":        m.Pair,
						},
						Payload: b,
					}
				default:
				}
			default:
			}
		}
	}()
	return c
}

type KrakenSubscriptionMessage struct {
	ChannelID    int    `json:"channelID"`
	ChannelName  string `json:"channelName"`
	Event        string `json:"event"`
	Pair         string `json:"pair"`
	Status       string `json:"status"`
	Subscription struct {
		Name string `json:"name"`
	} `json:"subscription"`
}

func parseKrakenSubscriptionMessage(b []byte) (KrakenSubscriptionMessage, error) {
	v := KrakenSubscriptionMessage{}
	err := json.Unmarshal(b, &v)
	return v, err
}

func parseKrakenTickerMessage(b []byte) (int, interface{}, error) {
	var v []interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return 0, nil, err
	}
	return int(v[0].(float64)), v[1], nil
}
