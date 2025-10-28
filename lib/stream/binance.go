package stream

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/jimtang2/bc-go/lib/db"
)

type BinanceStream struct {
	out    chan Message
	signal chan int
	socket *websocket.Conn
}

func NewBinanceStream() *BinanceStream {
	return &BinanceStream{
		out:    make(chan Message),
		signal: make(chan int),
	}
}

func (s *BinanceStream) Output() chan Message {
	return s.out
}

func (s *BinanceStream) Start() {
	go start(s)
	s.signal <- 1
}

func (s *BinanceStream) Name() string {
	return "binance"
}

func (s *BinanceStream) Signal() chan int {
	return s.signal
}

func (s *BinanceStream) Connect() error {
	return s.connect()
}

func (s *BinanceStream) Listen() {
	s.listen()
}

func (s *BinanceStream) Close() {
	if s.socket != nil {
		s.socket.Close()
	}
}

func (s *BinanceStream) connect() error {
	var err error
	s.socket, _, err = websocket.DefaultDialer.Dial(fmt.Sprintf("wss://stream.binance.com:9443/stream?streams=%s@ticker", strings.Join(db.TrackedPairs("binance"), "@ticker/")), nil)
	return err
}

func (s *BinanceStream) listen() {
	for {
		messageType, b, err := s.socket.ReadMessage()
		if err != nil {
			log.Println("[binance]", err)
			break
		}
		switch messageType {
		case websocket.PingMessage:
			err = s.socket.WriteMessage(websocket.PongMessage, b)
			if err != nil {
				log.Println("[binance]", err)
				break
			}
		case websocket.TextMessage:
			go s.send(b)
		default:
		}
	}
	s.signal <- 2
}

func (s *BinanceStream) send(b []byte, args ...interface{}) {
	v := &BinanceTicker{}
	if err := json.Unmarshal(b, v); err != nil {
		return
	}
	t, ok := v.Ticker().Fmt()
	if !ok {
		return
	}
	s.out <- Message{
		Topic:   "tickers",
		Key:     t.Key(),
		Payload: t.Bytes(),
	}
}

type BinanceTicker struct {
	Stream string `json:"stream"`
	Data   struct {
		EventType          string `json:"e"` // Event type
		EventTime          int64  `json:"E"` // Event time
		Symbol             string `json:"s"` // Symbol
		PriceChange        string `json:"p"` // Price change
		PriceChangePercent string `json:"P"` // Price change percent
		WeightedAverage    string `json:"w"` // Weighted average price
		PreviousClose      string `json:"x"` // Previous close
		LastPrice          string `json:"c"` // Last price
		LastQuantity       string `json:"Q"` // Last quantity
		BestBidPrice       string `json:"b"` // Best bid price
		BestBidQuantity    string `json:"B"` // Best bid quantity
		BestAskPrice       string `json:"a"` // Best ask price
		BestAskQuantity    string `json:"A"` // Best ask quantity
		OpenPrice          string `json:"o"` // Open price
		HighPrice          string `json:"h"` // High price
		LowPrice           string `json:"l"` // Low price
		Volume             string `json:"v"` // Volume
		QuoteVolume        string `json:"q"` // Quote volume
		OpenTime           int64  `json:"O"` // Statistics open time
		CloseTime          int64  `json:"C"` // Statistics close time
		FirstTradeID       int64  `json:"F"` // First trade ID
		LastTradeID        int64  `json:"L"` // Last trade ID
		TradeCount         int64  `json:"n"` // Total number of trades
	} `json:"data"`
}

func (v *BinanceTicker) Ticker() *Ticker {
	t := Ticker{}
	t.Exchange = "binance"
	t.Pair = v.Data.Symbol
	t.Bid, _ = strconv.ParseFloat(v.Data.BestBidPrice, 64)
	t.BidSize, _ = strconv.ParseFloat(v.Data.BestBidQuantity, 64)
	t.Ask, _ = strconv.ParseFloat(v.Data.BestAskPrice, 64)
	t.AskSize, _ = strconv.ParseFloat(v.Data.BestAskQuantity, 64)
	t.EventTime = v.Data.EventTime
	return &t
}

// BinanceTicker defines a detailed quote structure for centralized exchanges, including standard market data (https://developers.binance.com/docs/binance-spot-api-docs/web-socket-streams#individual-symbol-ticker-streams)
/*{
  "e": "24hrTicker",  // Event type
  "E": 1672515782136, // Event time
  "s": "BNBBTC",      // Symbol
  "p": "0.0015",      // Price change
  "P": "250.00",      // Price change percent
  "w": "0.0018",      // Weighted average price
  "x": "0.0009",      // First trade(F)-1 price (first trade before the 24hr rolling window)
  "c": "0.0025",      // Last price
  "Q": "10",          // Last quantity
  "b": "0.0024",      // Best bid price
  "B": "10",          // Best bid quantity
  "a": "0.0026",      // Best ask price
  "A": "100",         // Best ask quantity
  "o": "0.0010",      // Open price
  "h": "0.0025",      // High price
  "l": "0.0010",      // Low price
  "v": "10000",       // Total traded base asset volume
  "q": "18",          // Total traded quote asset volume
  "O": 0,             // Statistics open time
  "C": 86400000,      // Statistics close time
  "F": 0,             // First trade ID
  "L": 18150,         // Last trade Id
  "n": 18151          // Total number of trades
}*/
