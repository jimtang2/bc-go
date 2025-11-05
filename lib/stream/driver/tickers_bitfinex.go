package driver

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jimtang2/bc-go/lib/kafka"
	"github.com/jimtang2/bc-go/lib/stream"
)

func init() {
	stream.Register("bitfinex-tickers", &BitfinexTickers{
		channels: map[int]string{},
		mu:       sync.Mutex{},
	})
}

type BitfinexTickers struct {
	ws       *websocket.Conn
	channels map[int]string // map of channel id to pair text
	mu       sync.Mutex
}

func (s *BitfinexTickers) Open(pairs []string) (*websocket.Conn, error) {
	ws, _, err := websocket.DefaultDialer.Dial("wss://api-pub.bitfinex.com/ws/2", nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	s.ws = ws
	m := map[string]interface{}{
		"event":   "subscribe",
		"channel": "ticker",
	}
	for _, pair := range pairs {
		m["symbol"] = "t" + pair
		if err := ws.WriteJSON(m); err != nil {
			log.Println(err)
		}
	}
	return ws, nil
}

func (s *BitfinexTickers) Close() error {
	if s.ws == nil {
		return nil
	} else {
		return s.ws.Close()
	}
}

func (s *BitfinexTickers) OnWebsocketMessage(messageType int, b []byte) (kafka.KMessage, error) {
	if messageType == websocket.TextMessage {
		if len(b) == 0 {
			return nil, nil
		} else if b[0] == '[' {
			v := []interface{}{}
			if err := json.Unmarshal(b, &v); err != nil || len(v) < 2 {
				return nil, err
			}
			chanID, ok := v[0].(float64)
			if !ok {
				return nil, fmt.Errorf("unexpected data type")
			}
			data, ok := v[1].([]interface{})
			if !ok {
				// api returned heartbeat; ignoring
				return nil, nil
			} else if len(data) < 10 {
				return nil, fmt.Errorf("[bitfinex] unexpected number of data points %v %v", len(b), string(b))
			}
			vv := &BitfinexTicker{
				ChanID:              int(chanID),
				Bid:                 data[0].(float64),
				BidSize:             data[1].(float64),
				Ask:                 data[2].(float64),
				AskSize:             data[3].(float64),
				DailyChange:         data[4].(float64),
				DailyChangeRelative: data[5].(float64),
				LastPrice:           data[6].(float64),
				Volume:              data[7].(float64),
				High:                data[8].(float64),
				Low:                 data[9].(float64),
			}
			t := vv.TickerMessage()
			t.Pair = s.channels[int(chanID)]
			t.Fmt()
			return t, nil
		} else if b[0] == '{' {
			v := BitfinexSubscribeAck{}
			if err := json.Unmarshal(b, &v); err != nil {
				log.Println(err)
				return nil, nil
			}
			if v.Event == "subscribed" {
				s.mu.Lock()
				s.channels[v.ChanID] = v.Pair
				s.mu.Unlock()
				// log.Println("[bitfinex] subscribed to", v.Pair)
			} else if v.Event == "info" {
				// log.Println(string(b))
				// 20051: need restart websocket connection
				// 20060: maintenance start
				// 20061: maintenance end
				// log.Printf("[bitfinex] [%v] %v", m.Code, string(b))
			}
		}
	}
	return nil, nil
}

type BitfinexSubscribeAck struct {
	Event   string `json:"event"`
	ChanID  int    `json:"chanId"`
	Pair    string `json:"pair"`
	Channel string `json:"-"`
	Symbol  string `json:"-"`
	Msg     string `json:"-"`
	Code    int    `json:"-"`
}

type BitfinexTicker struct {
	ChanID              int
	Bid                 float64
	BidSize             float64
	Ask                 float64
	AskSize             float64
	DailyChange         float64
	DailyChangeRelative float64
	LastPrice           float64
	Volume              float64
	High                float64
	Low                 float64
}

func (v *BitfinexTicker) TickerMessage() *TickerMessage {
	t := TickerMessage{
		Exchange: "bitfinex",
	}
	t.Bid = v.Bid
	t.BidSize = v.BidSize
	t.Ask = v.Ask
	t.AskSize = v.AskSize
	return &t
}

// https://docs.bitfinex.com/docs/ws-public
// https://docs.bitfinex.com/reference/ws-public-ticker
// https://docs.bitfinex.com/docs/ws-general#subscribe-to-channels
// bitfinex websocket sends as response arrays of items that either number, string or array of number; e.g.: [1234,"hb"] or [1234,[1,2,3,4,5,6,7,8,9,10]]
// this function returns the channel id integer, a boolean indicating whether message is heartbeat, and error
/*
Index 	Field 	Type 	Description
[0]	BID	Float	Price of last highest bid
[1]	BID_SIZE	Float	Sum of the 25 highest bid sizes
[2]	ASK	Float	Price of last lowest ask
[3]	ASK_SIZE	Float	Sum of the 25 lowest ask sizes
[4]	DAILY_CHANGE	Float	Amount that the last price has changed since yesterday
[5]	DAILY_CHANGE_RELATIVE	Float	Relative price change since yesterday (*100 for percentage change)
[6]	LAST_PRICE	Float	Price of the last trade.
[7]	VOLUME	Float	Daily volume
[8]	HIGH	Float	Daily high
[9]	LOW	Float	Daily low
*/
