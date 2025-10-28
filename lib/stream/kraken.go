package stream

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/jimtang2/bc-go/lib/db"
)

type KrakenStream struct {
	messages chan Message
}

func (s *KrakenStream) Start() chan Message {
	s.messages = make(chan Message)
	ws, err := s.subscribe()
	if err != nil {
		log.Println("[kraken]", err)
		return s.messages
	}
	go func() {
		defer ws.Close()
		for {
			_, b, err := ws.ReadMessage()
			if err != nil {
				log.Println("[kraken]", err)
				return
			}
			go s.send(b)
		}
	}()
	return s.messages
}

func (s *KrakenStream) subscribe() (*websocket.Conn, error) {
	ws, _, err := websocket.DefaultDialer.Dial("wss://ws.kraken.com/v2", nil)
	m := struct {
		Method string `json:"method"`
		Params struct {
			Channel string   `json:"channel"`
			Symbol  []string `json:"symbol"`
		} `json:"params"`
	}{}
	if err != nil {
		return nil, err
	}
	m.Method = "subscribe"
	m.Params.Channel = "ticker"
	m.Params.Symbol = db.TrackedPairs("kraken")
	return ws, ws.WriteJSON(m)
}

func (s *KrakenStream) send(b []byte, args ...interface{}) {
	v := KrakenTicker{}
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

// https://docs.kraken.com/api/docs/websocket-v2/ticker
/*{
  "channel": "ticker",
  "type": "snapshot",
  "data": [{
    "symbol": "ALGO/USD",
    "bid": 0.10025,
    "bid_qty": 740.0,
    "ask": 0.10036,
    "ask_qty": 1361.44813783,
    "last": 0.10035,
    "volume": 997038.98383185,
    "vwap": 0.10148,
    "low": 0.09979,
    "high": 0.10285,
    "change": -0.00017,
    "change_pct": -0.17
  }]
}*/
type KrakenTicker struct {
	Channel string `json:"channel"`
	Type    string `json:"type"`
	Data    []struct {
		Symbol        string  `json:"symbol"`
		Bid           float64 `json:"bid"`
		BidSize       float64 `json:"bid_qty"`
		Ask           float64 `json:"ask"`
		AskSize       float64 `json:"ask_qty"`
		Last          float64 `json:"last"`
		Volume        float64 `json:"volume"`
		Vwap          float64 `json:"vwap"`
		Low           float64 `json:"low"`
		High          float64 `json:"high"`
		Change        float64 `json:"change"`
		ChangePercent float64 `json:"change_pct"`
	} `json:"data"`
}

func (v *KrakenTicker) Ticker() *Ticker {
	t := &Ticker{
		Exchange: "kraken",
	}
	if len(v.Data) < 1 {
		return t
	}
	t.Pair = v.Data[0].Symbol
	t.Bid = v.Data[0].Bid
	t.BidSize = v.Data[0].BidSize
	t.Ask = v.Data[0].Ask
	t.AskSize = v.Data[0].AskSize
	return t
}

/*
	{
		"method": "subscribe",
		"params": {
	    "channel": "ticker",
	    "symbol": [
	      "ALGO/USD"
	    ]
		}
	}
*/

// type KrakenEventMessage struct {
// 	ChannelID    int    `json:"channelID"`
// 	ChannelName  string `json:"channelName"`
// 	Event        string `json:"event"`
// 	Pair         string `json:"pair"`
// 	Status       string `json:"status"`
// 	Subscription struct {
// 		Name string `json:"name"`
// 	} `json:"subscription"`
// }

// func (p *KrakenStream) Pair(b []byte) string {
// 	v := struct {
// 		Channel string `json:"channel"`
// 		Type    string `json:"type"`
// 		Method  string `json:"method"`
// 		Result  struct {
// 			Symbol string `json:"symbol"`
// 		} `json:"result"`
// 		Success bool `json:"success"`
// 		Data    []struct {
// 			Symbol string `json:"symbol"`
// 		} `json:"data"`
// 	}{}
// 	err := json.Unmarshal(b, &v)
// 	if err != nil {
// 		log.Println("[kraken]", err)
// 		return ""
// 	}
// 	if v.Method == "subscribe" {
// 		log.Println("[kraken] subscribed to", v.Result.Symbol)
// 		return ""
// 	} else if v.Channel == "ticker" {
// 		return v.Data[0].Symbol
// 	} else {
// 		return ""
// 	}
// }
