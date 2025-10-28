package stream

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jimtang2/bc-go/lib/db"
)

type CoinbaseStream struct {
	messages chan Message
}

func (s *CoinbaseStream) Start() chan Message {
	s.messages = make(chan Message)
	ws, err := s.subscribe()
	if err != nil {
		log.Println("[coinbase]", err)
		return s.messages
	}
	go func() {
		defer ws.Close()
		for {
			messageType, b, err := ws.ReadMessage()
			if err != nil {
				log.Println("[coinbase]", err)
				return
			}
			switch messageType {
			case websocket.TextMessage:
				go s.send(b)
			default:
			}
		}
	}()
	return s.messages
}

func (s *CoinbaseStream) subscribe() (*websocket.Conn, error) {
	ws, _, err := websocket.DefaultDialer.Dial("wss://ws-feed.exchange.coinbase.com", nil)
	if err != nil {
		return nil, err
	}
	subscribeMsg := map[string]interface{}{
		"type":        "subscribe",
		"product_ids": db.TrackedPairs("coinbase"),
		"channels":    []string{"ticker"},
	}
	if err := ws.WriteJSON(subscribeMsg); err != nil {
		return nil, err
	}
	return ws, nil
}

func (s *CoinbaseStream) send(b []byte, args ...interface{}) {
	v := CoinbaseTicker{}
	if err := json.Unmarshal(b, &v); err != nil {
		return
	}
	t, ok := v.Ticker().Fmt()
	if !ok {
		return
	}
	s.messages <- Message{
		Topic:   "tickers",
		Key:     t.Key(),
		Payload: t.Bytes(),
	}
}

// https://docs.cdp.coinbase.com/exchange/websocket-feed/channels#ticker-channel
func (p *CoinbaseStream) pair(b []byte) (string, bool) {
	v := struct {
		ProductID string `json:"product_id"`
	}{}
	err := json.Unmarshal(b, &v)
	return v.ProductID, err == nil
}

/*
	{
	  "type": "ticker",
	  "sequence": 37475248783,
	  "product_id": "ETH-USD",
	  "price": "1285.22",
	  "open_24h": "1310.79",
	  "volume_24h": "245532.79269678",
	  "low_24h": "1280.52",
	  "high_24h": "1313.8",
	  "volume_30d": "9788783.60117027",
	  "best_bid": "1285.04",
	  "best_bid_size": "0.46688654",
	  "best_ask": "1285.27",
	  "best_ask_size": "1.56637040",
	  "side": "buy",
	  "time": "2022-10-19T23:28:22.061769Z",
	  "trade_id": 370843401,
	  "last_size": "11.4396987"
	}
*/
type CoinbaseTicker struct {
	Type        string `json:"type"`
	Sequence    int    `json:"sequence"`
	ProductID   string `json:"product_id"`
	Price       string `json:"price"`
	Open24h     string `json:"open_24h"`
	Volume24h   string `json:"volume_24h"`
	Low24h      string `json:"low_24h"`
	High24h     string `json:"high_24h"`
	Volume30d   string `json:"volume_30d"`
	BestBid     string `json:"best_bid"`
	BestBidSize string `json:"best_bid_size"`
	BestAsk     string `json:"best_ask"`
	BestAskSize string `json:"best_ask_size"`
	Side        string `json:"side"`
	Time        string `json:"time"`
	TradeID     int    `json:"trade_id"`
	LastSize    string `json:"last_size"`
}

func (v *CoinbaseTicker) Ticker() *Ticker {
	t := Ticker{
		Exchange: "coinbase",
	}
	t.Pair = v.ProductID
	t.Bid, _ = strconv.ParseFloat(v.BestBid, 64)
	t.BidSize, _ = strconv.ParseFloat(v.BestBidSize, 64)
	t.Ask, _ = strconv.ParseFloat(v.BestAsk, 64)
	t.AskSize, _ = strconv.ParseFloat(v.BestAskSize, 64)
	ts, err := time.Parse(time.RFC3339Nano, v.Time)
	if err == nil {
		t.EventTime = ts.UnixMilli()
	}
	return &t
}
