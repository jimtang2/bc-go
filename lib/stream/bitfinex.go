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

type BitfinexStream struct {
	// bitfinex uses an id in payloads instead of the pair text; this requires caching the subscription event message
	subscriptions map[int]BitfinexEventMessage
	mu            sync.Mutex
	parser        *BitfinexParser
}

// https://docs.bitfinex.com/docs/ws-general#subscribe-to-channels
func (s *BitfinexStream) Start() chan Message {
	s.subscriptions = map[int]BitfinexEventMessage{}
	s.mu = sync.Mutex{}
	c := make(chan Message)
	pairs := db.TrackedPairs("bitfinex")
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
	for _, pair := range pairs {
		subscribeMsg["symbol"] = "t" + pair
		if err := ws.WriteJSON(subscribeMsg); err != nil {
			log.Println("[bitfinex]", err)
		}
	}
	log.Println("[bitfinex] subscribed to", strings.Join(pairs, ", "))
	go func() {
		defer ws.Close()
		for {
			_, b, err := ws.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			if len(b) == 0 {
				log.Println("[bitfinex] unexpected empty message")
				continue
			}
			firstChar := string(b[0])
			if firstChar == "{" {
				m, err := s.parser.EventMessage(b)
				if err != nil {
					log.Println(err)
					continue
				}
				if m.Event == "subscribed" {
					s.mu.Lock()
					s.subscriptions[m.ChanID] = m
					s.mu.Unlock()
					log.Println("[bitfinex] subscribed to", m.Pair)
				} else if m.Event == "info" {
					// 20051: need restart websocket connection
					// 20060: maintenance start
					// 20061: maintenance end
					log.Printf("[bitfinex] [%v] %v", m.Code, string(b))
				}
			} else if firstChar == "[" {
				chanID, isHeartbeat, err := s.parser.ChanID(b)
				if err != nil {
					log.Println("[bitfinex]", err)
					continue
				} else if isHeartbeat {
					continue
				}
				s.mu.Lock()
				m, ok := s.subscriptions[chanID]
				s.mu.Unlock()
				if !ok {
					log.Println("[bitfinex] unexpected chanID not set:", string(b))
					continue
				}
				go func() {
					c <- Message{
						Topic: "tickers_raw",
						Key:   "[bitfinex]" + m.Pair,
						Headers: map[string]string{
							"exchange": "bitfinex",
							"pair":     m.Pair,
						},
						Payload: b,
					}
				}()
				go func() {
					v, err := (&BitfinexParser{}).Ticker(b)
					if err != nil {
						log.Println(err)
						return
					}
					t := v.Ticker().Fmt()
					t.Pair = m.Pair
					if !t.IsValid() {
						return
					}
					c <- Message{
						Topic:   "tickers",
						Key:     t.Key(),
						Payload: t.Bytes(),
					}
				}()
			}
		}
	}()
	return c
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
	ChanID              float64
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

type BitfinexEventMessage struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	ChanID  int    `json:"chanId"`
	Symbol  string `json:"symbol"`
	Pair    string `json:"pair"`
	Msg     string `json:"msg"`
	Code    int    `json:"code"`
}

type BitfinexParser struct{}

func (p *BitfinexParser) EventMessage(b []byte) (BitfinexEventMessage, error) {
	v := BitfinexEventMessage{}
	err := json.Unmarshal(b, &v)
	return v, err
}

// https://docs.bitfinex.com/reference/ws-public-ticker
// bitfinex websocket sends as response arrays of items that either number, string or array of number; e.g.: [1234,"hb"] or [1234,[1,2,3,4,5,6,7,8,9,10]]
// this function returns the channel id integer, a boolean indicating whether message is heartbeat, and error
func (p *BitfinexParser) ChanID(b []byte) (int, bool, error) {
	var v []interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return 0, false, err
	}
	_, isHeartbeat := v[1].(string)
	return int(v[0].(float64)), isHeartbeat, nil
}

func (p *BitfinexParser) Ticker(b []byte) (*BitfinexTicker, error) {
	v := []interface{}{}
	if err := json.Unmarshal(b, &v); err != nil {
		return nil, err
	}
	if len(v) < 2 {
		return nil, fmt.Errorf("unexpected data length:", len(v))
	}
	vv, ok := v[1].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data content", string(b))
	}
	return &BitfinexTicker{
		ChanID:              v[0].(float64),
		Bid:                 vv[0].(float64),
		BidSize:             vv[1].(float64),
		Ask:                 vv[2].(float64),
		AskSize:             vv[3].(float64),
		DailyChange:         vv[4].(float64),
		DailyChangeRelative: vv[5].(float64),
		LastPrice:           vv[6].(float64),
		Volume:              vv[7].(float64),
		High:                vv[8].(float64),
		Low:                 vv[9].(float64),
	}, nil
}
