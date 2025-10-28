package stream

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jimtang2/bc-go/lib/db"
)

type BitfinexStream struct {
	channels map[int]string // map of channel id to pair text
	mu       sync.Mutex
	messages chan Message
}

// https://docs.bitfinex.com/docs/ws-general#subscribe-to-channels
func (s *BitfinexStream) Start() chan Message {
	s.channels = map[int]string{}
	s.mu = sync.Mutex{}
	s.messages = make(chan Message)
	ws, err := s.subscribe()
	if err != nil {
		log.Println("[bitfinex]", err)
		return s.messages
	}
	go func() {
		defer ws.Close()
		for {
			messageType, b, err := ws.ReadMessage()
			if err != nil {
				log.Println("[bitfinex]", err)
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

// https://docs.bitfinex.com/docs/ws-public
// https://docs.bitfinex.com/reference/ws-public-ticker
// bitfinex websocket sends as response arrays of items that either number, string or array of number; e.g.: [1234,"hb"] or [1234,[1,2,3,4,5,6,7,8,9,10]]
// this function returns the channel id integer, a boolean indicating whether message is heartbeat, and error
func (s *BitfinexStream) subscribe() (*websocket.Conn, error) {
	ws, _, err := websocket.DefaultDialer.Dial("wss://api-pub.bitfinex.com/ws/2", nil)
	if err != nil {
		return nil, err
	}
	subscribeMsg := map[string]interface{}{
		"event":   "subscribe",
		"channel": "ticker",
	}
	for _, pair := range db.TrackedPairs("bitfinex") {
		subscribeMsg["symbol"] = "t" + pair
		if err := ws.WriteJSON(subscribeMsg); err != nil {
			log.Println("[bitfinex]", err)
		}
	}
	return ws, nil
}

func (s *BitfinexStream) send(b []byte, args ...interface{}) {
	if len(b) == 0 {
		return
	}
	switch b[0] {
	case '{':
		v := struct {
			Event   string `json:"event"`
			ChanID  int    `json:"chanId"`
			Pair    string `json:"pair"`
			Channel string `json:"-"`
			Symbol  string `json:"-"`
			Msg     string `json:"-"`
			Code    int    `json:"-"`
		}{}
		if err := json.Unmarshal(b, &v); err != nil {
			return
		}
		if v.Event == "subscribed" {
			s.mu.Lock()
			s.channels[v.ChanID] = v.Pair
			s.mu.Unlock()
			// log.Println("[bitfinex] subscribed to", m.Pair)
		} else if v.Event == "info" {
			// 20051: need restart websocket connection
			// 20060: maintenance start
			// 20061: maintenance end
			// log.Printf("[bitfinex] [%v] %v", m.Code, string(b))
		}
	case '[':
		v := []interface{}{}
		if err := json.Unmarshal(b, &v); err != nil || len(v) < 2 {
			return
		}
		chanID, ok := v[0].(float64)
		if !ok {
			return
		}
		data, ok := v[1].([]interface{})
		if !ok || len(data) < 10 {
			return
		}
		t := (&BitfinexTicker{
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
		}).Ticker()
		t.Pair = s.channels[int(chanID)]
		if t, ok = t.Fmt(); !ok {
			return
		}
		s.messages <- Message{
			Topic:   "tickers",
			Key:     t.Key(),
			Payload: t.Bytes(),
			Headers: map[string]string{
				"exchange": "bitfinex",
			},
		}
	default:
	}
}

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

func (v *BitfinexTicker) Ticker() *Ticker {
	t := Ticker{
		Exchange: "bitfinex",
	}
	t.Bid = v.Bid
	t.BidSize = v.BidSize
	t.Ask = v.Ask
	t.AskSize = v.AskSize
	return &t
}
