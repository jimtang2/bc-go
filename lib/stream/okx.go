package stream

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jimtang2/bc-go/lib/db"
)

type OKX struct {
	lastMessageTime time.Time
	mu              sync.Mutex
}

func (s *OKX) Start() chan Message {
	c := make(chan Message)
	url := "wss://ws.okx.com:8443/ws/v5/public"
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("[okx]", err)
		return c
	}
	if err := OKXSubscribe(ws); err != nil {
		log.Println("[okx]", err)
		return c
	}
	// OKX handles ping-pong if the connection is idle (no new message) after 30s
	// it expects a message with a content "ping" (different than ping message) and would respond with the same with "pong", to determine that the connection is stale or needs reconnect
	// the following goroutine doesn't accomplish much but doesn't hurt to leave
	// if/when basic OKX stream behavior requires more work this mechanism may come to use
	go func() {
		for {
			s.mu.Lock()
			t := s.lastMessageTime
			s.mu.Unlock()
			if time.Now().Sub(t) > 30*time.Second {
				if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					log.Println("[okx]", err)
					return
				}
			}
			time.Sleep(2 * time.Second)
		}
	}()
	//
	go func() {
		defer ws.Close()
		for {
			messageType, b, err := ws.ReadMessage()
			if err != nil {
				log.Println("[okx]", err)
				return
			}
			s.mu.Lock()
			s.lastMessageTime = time.Now()
			s.mu.Unlock()
			switch messageType {
			case websocket.PongMessage:
				continue
			case websocket.TextMessage:
				m, err := OKXParseResponse(b)
				if err != nil {
					log.Println(err)
					continue
				}
				if m.Event == "subscribe" {
					log.Println("[okx] subscribed to", m.Arg.InstID)
				} else {
					if m.Arg.Channel == "tickers" {
						c <- Message{
							Topic: "tickers",
							Key:   fmt.Sprintf("[okx]%s", m.Arg.InstID),
							Headers: map[string]string{
								"exchange": "okx",
								"api":      "spot order book market data tickers channel",
								"pair":     m.Arg.InstID,
							},
							Payload: b,
						}
					}
				}
			default:
			}
		}
	}()

	return c
}

/*
example payload:

	{
	  "arg": {
	    "channel": "tickers",
	    "instId": "BTC-USDT"
	  },
	  "data": [
	    {
	      "instType": "SPOT",
	      "instId": "BTC-USDT",
	      "last": "9999.99",
	      "lastSz": "0.1",
	      "askPx": "9999.99",
	      "askSz": "11",
	      "bidPx": "8888.88",
	      "bidSz": "5",
	      "open24h": "9000",
	      "high24h": "10000",
	      "low24h": "8888.88",
	      "volCcy24h": "2222",
	      "vol24h": "2222",
	      "sodUtc0": "2222",
	      "sodUtc8": "2222",
	      "ts": "1597026383085"
	    }
	  ]
	}

payload description:
arg 	Object 	Successfully subscribed channel
> channel 	String 	Channel name
> instId 	String 	Instrument ID
data 	Array of objects 	Subscribed data
> instType 	String 	Instrument type
> instId 	String 	Instrument ID
> last 	String 	Last traded price
> lastSz 	String 	Last traded size. 0 represents there is no trading volume
> askPx 	String 	Best ask price
> askSz 	String 	Best ask size
> bidPx 	String 	Best bid price
> bidSz 	String 	Best bid size
> open24h 	String 	Open price in the past 24 hours
> high24h 	String 	Highest price in the past 24 hours
> low24h 	String 	Lowest price in the past 24 hours
> volCcy24h 	String 	24h trading volume, with a unit of currency.
If it is a derivatives contract, the value is the number of base currency.
If it is SPOT/MARGIN, the value is the quantity in quote currency.
> vol24h 	String 	24h trading volume, with a unit of contract.
If it is a derivatives contract, the value is the number of contracts.
If it is SPOT/MARGIN, the value is the quantity in base currency.
> sodUtc0 	String 	Open price in the UTC 0
> sodUtc8 	String 	Open price in the UTC 8
> ts 	String 	Ticker data generation time, Unix timestamp format in milliseconds, e.g. 1597026383085
*/
type OKXResponse struct {
	Event string `json:"event"` // subscribe, unsubscribe, error
	Arg   struct {
		Channel string `json:"channel"`
		InstID  string `json:"instId"`
	} `json:"arg"`
	Code   string        `json:"code"`
	Msg    string        `json:"msg"`
	ConnID string        `json:"connId"`
	Data   OKXTickerData `json:"-"`
}

type OKXTickerData struct {
	InstType  string `json:"instType"`
	InstID    string `json:"instId"`
	Last      string `json:"last"`
	LastSz    string `json:"lastSz"`
	AskPx     string `json:"askPx"`
	AskSz     string `json:"askSz"`
	BidPx     string `json:"bidPx"`
	BidSz     string `json:"bidSz"`
	Open24h   string `json:"open24h"`
	High24h   string `json:"high24h"`
	Low24h    string `json:"low24h"`
	VolCcy24h string `json:"volCcy24h"`
	Vol24h    string `json:"vol24h"`
	SodUtc0   string `json:"sodUtc0"`
	SodUtc8   string `json:"sodUtc8"`
	Ts        string `json:"ts"`
}

type OKXSubscriptionRequest struct {
	Op   string        `json:"op"`
	Args []interface{} `json:"args"`
}

// OKX requires subscription in JSON format with op: "subscribe"
// https://www.okx.com/docs-v5/en/?language=shell#order-book-trading-market-data-ws-tickers-channel
func OKXSubscribe(ws *websocket.Conn) error {
	pairs := db.TrackedPairs("okx")
	m := OKXSubscriptionRequest{
		Op: "subscribe",
	}
	args := []interface{}{}
	for _, pair := range pairs {
		args = append(args, map[string]interface{}{
			"channel":  "tickers",
			"instId":   pair,
			"instType": "SPOT",
		})
	}
	m.Args = args
	if err := ws.WriteJSON(m); err != nil {
		return err
	}
	return nil
}

func OKXParseResponse(b []byte) (OKXResponse, error) {
	v := OKXResponse{}
	err := json.Unmarshal(b, &v)
	return v, err
}
