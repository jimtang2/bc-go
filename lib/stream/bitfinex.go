package stream

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jimtang2/bc-go/lib/db"
)

type Bitfinex struct {
	pairs map[int]BitfinexSubscriptionMessage
	mu    sync.Mutex
}

// https://docs.bitfinex.com/docs/ws-general#subscribe-to-channels
func (s *Bitfinex) Start() chan Message {
	s.pairs = map[int]BitfinexSubscriptionMessage{}
	s.mu = sync.Mutex{}
	c := make(chan Message)
	trackedPairs := db.TrackedPairs("bitfinex")
	url := "wss://api-pub.bitfinex.com/ws/2"
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("[bitfinex]", err)
		return c
	}
	// https://docs.bitfinex.com/docs/ws-public
	subscribeMsg := map[string]interface{}{
		"event":   "subscribe",
		"channel": "ticker",
	}
	for _, pair := range trackedPairs {
		subscribeMsg["symbol"] = "t" + pair
		if err := ws.WriteJSON(subscribeMsg); err != nil {
			log.Println("[bitfinex]", err)
		}
	}
	log.Println("[bitfinex] subscribed to", strings.Join(trackedPairs, ", "))
	go func() {
		defer ws.Close()
		for {
			messageType, b, err := ws.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			switch messageType {
			// bitfinex uses only one message type
			case websocket.TextMessage:
				if len(b) == 0 {
					log.Println("[bitfinex] unexpected empty message")
					continue
				}
				switch string(b[0]) {
				// json object type means subscription message
				case "{":
					m, err := parseBitfinexSubscriptionMessage(b)
					if err != nil {
						log.Println(err)
						continue
					}
					switch m.Event {
					case "subscribed":
						s.mu.Lock()
						s.pairs[m.ChanID] = m
						s.mu.Unlock()
						log.Println("[bitfinex] subscribed to", m.Pair)
					case "info":
						switch m.Code {
						case 20051: // need restart websocket connection
						case 20060: // maintenance start
						case 20061: // maintenance end
						default:
							log.Println("[bitfinex]", string(b))
						}

					default:
					}
				// json array type can be ticker, heartbeat
				case "[":
					pairID, val, err := parseBitfinexTickerMessage(b)
					if err != nil {
						log.Println(err)
						continue
					}
					if _, ok := val.(string); ok {
						continue
					}
					s.mu.Lock()
					p, ok := s.pairs[pairID]
					s.mu.Unlock()
					if !ok {
						continue
					}
					c <- Message{
						Topic: "tickers",
						Key:   "[bitfinex]" + p.Pair,
						Headers: map[string]string{
							"exchange": "bitfinex",
							"api":      "ticker stream",
							"chanId":   fmt.Sprintf(`%v`, p.ChanID),
							"channel":  p.Channel,
							"symbol":   p.Symbol,
							"pair":     p.Pair,
							"event":    p.Event,
						},
						Payload: b,
					}
				default:

				}
			case websocket.PingMessage:
			case websocket.CloseMessage:
			default:
			}
		}
	}()
	return c
}

type BitfinexSubscriptionMessage struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	ChanID  int    `json:"chanId"`
	Symbol  string `json:"symbol"`
	Pair    string `json:"pair"`
	Msg     string `json:"msg"`
	Code    int    `json:"code"`
}

func parseBitfinexSubscriptionMessage(b []byte) (BitfinexSubscriptionMessage, error) {
	v := BitfinexSubscriptionMessage{}
	err := json.Unmarshal(b, &v)
	return v, err
}

// https://docs.bitfinex.com/reference/ws-public-ticker
func parseBitfinexTickerMessage(b []byte) (int, interface{}, error) {
	var v []interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return 0, nil, err
	}
	return int(v[0].(float64)), v[1], nil
}
