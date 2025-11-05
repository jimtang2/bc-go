package driver

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/jimtang2/bc-go/lib/kafka"
	"github.com/jimtang2/bc-go/lib/stream"
)

func init() {
	stream.Register("okx-tickers", &OKXTickers{})
}

type OKXTickers struct {
	ws *websocket.Conn
}

func (s *OKXTickers) Close() error {
	if s.ws != nil {
		return s.ws.Close()
	} else {
		return nil
	}
}

func (s *OKXTickers) Open(pairs []string) (*websocket.Conn, error) {
	ws, _, err := websocket.DefaultDialer.Dial("wss://ws.okx.com:8443/ws/v5/public", nil)
	if err != nil {
		return nil, err
	}
	args := []interface{}{}
	for _, pair := range pairs {
		args = append(args, map[string]interface{}{
			"channel":  "tickers",
			"instId":   pair,
			"instType": "SPOT",
		})
	}
	m := map[string]interface{}{
		"op":   "subscribe",
		"args": interface{}(args),
	}
	return ws, ws.WriteJSON(m)
}

func (s *OKXTickers) OnWebsocketMessage(messageType int, b []byte) (kafka.KMessage, error) {
	v := OKXTicker{}
	if err := json.Unmarshal(b, &v); err != nil {
		log.Println("[okx]", err)
		return nil, nil
	}
	if v.Event == "subscribe" {
		return nil, nil
	} else if v.Arg.Channel == "tickers" {
		return v.TickerMessage(), nil
	}
	return nil, nil
}

type OKXTicker struct {
	Event  string `json:"event"`  // from event response, missing in ticker; values are subscribe, unsubscribe, error
	Code   string `json:"code"`   // from event response, missing in ticker
	Msg    string `json:"msg"`    // from event response, missing in ticker
	ConnID string `json:"connId"` // from event response, missing in ticker
	Arg    struct {
		Channel string `json:"channel"`
		InstID  string `json:"instId"`
	} `json:"arg"`
	Data []struct {
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
	} `json:"data"`
}

func (v *OKXTicker) TickerMessage() *TickerMessage {
	t := &TickerMessage{
		Exchange: "okx",
	}
	t.Pair = v.Arg.InstID
	t.Bid, _ = strconv.ParseFloat(v.Data[0].BidPx, 64)
	t.BidSize, _ = strconv.ParseFloat(v.Data[0].BidSz, 64)
	t.Ask, _ = strconv.ParseFloat(v.Data[0].AskPx, 64)
	t.AskSize, _ = strconv.ParseFloat(v.Data[0].AskSz, 64)
	t.EventTime, _ = strconv.ParseInt(v.Data[0].Ts, 10, 64)
	t.Fmt()
	return t
}

// OKX requires subscription in JSON format with op: "subscribe"
// https://www.okx.com/docs-v5/en/?language=shell#order-book-trading-market-data-ws-tickers-channel
/*
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
/*{
  "arg": {
    "channel": "tickers",
    "instId": "BTC-USDT"
  },
  "data": [{
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
  }]
}*/
