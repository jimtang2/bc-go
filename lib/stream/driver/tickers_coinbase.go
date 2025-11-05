package driver

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jimtang2/bc-go/lib/kafka"
	"github.com/jimtang2/bc-go/lib/stream"
	"github.com/jimtang2/bc-go/pkg/pb/v1"
)

func init() {
	stream.Register("coinbase-tickers", &CoinbaseTickers{})
}

type CoinbaseTickers struct {
	ws *websocket.Conn
}

func (s *CoinbaseTickers) Open(pairs []string) (*websocket.Conn, error) {
	ws, _, err := websocket.DefaultDialer.Dial("wss://ws-feed.exchange.coinbase.com", nil)
	if err != nil {
		return nil, err
	}
	s.ws = ws
	err = ws.WriteJSON(map[string]interface{}{
		"type":        "subscribe",
		"product_ids": pairs,
		"channels":    []string{"ticker"},
	})
	return ws, err
}

func (s *CoinbaseTickers) Close() error {
	if s.ws == nil {
		return nil
	} else {
		return s.ws.Close()
	}
}

func (s *CoinbaseTickers) OnWebsocketMessage(messageType int, b []byte) (kafka.KMessage, error) {
	if messageType == websocket.TextMessage {
		v := CoinbaseTicker{}
		if err := json.Unmarshal(b, &v); err != nil {
			log.Println(err)
			return nil, nil
		}
		return v.Ticker(), nil
	} else {
		return nil, nil
	}
}

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
	t := &Ticker{
		Ticker: &pb.Ticker{
			Exchange: "coinbase",
			Pair:     v.ProductID,
		},
	}
	t.Ticker.Bid, _ = strconv.ParseFloat(v.BestBid, 64)
	t.Ticker.BidSize, _ = strconv.ParseFloat(v.BestBidSize, 64)
	t.Ticker.Ask, _ = strconv.ParseFloat(v.BestAsk, 64)
	t.Ticker.AskSize, _ = strconv.ParseFloat(v.BestAskSize, 64)
	ts, err := time.Parse(time.RFC3339Nano, v.Time)
	if err == nil {
		t.Ticker.EventTime = ts.UnixMilli()
	}
	t.Fmt()
	return t
}

// https://docs.cdp.coinbase.com/exchange/websocket-feed/channels#ticker-channel
/*{
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
}*/
