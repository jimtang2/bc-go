package stream

import (
	"encoding/json"
	"log"
	"strings"
	"time"
)

type Message struct {
	Topic   string
	Key     string
	Headers map[string]string
	Payload []byte
}

const (
	SIGNAL_EXIT = iota
	SIGNAL_START
	SIGNAL_RESTART
)

type Stream interface {
	Start()               // called by Proxy.Stream()
	Output() chan Message // called by Proxy.Stream()
	Name() string         // called by stream.Start()
	Signal() chan int     // called by stream.Start()
	Connect() error       // called by stream.Start()
	Listen()              // called by stream.Start()
	Close()               // called by stream.Start()
}

func start(s Stream) {
	for {
		switch <-s.Signal() {
		case SIGNAL_EXIT:
			s.Close()
			log.Printf("[%v] connection closed (0)", s.Name())
			return
		case SIGNAL_START:
			if err := s.Connect(); err != nil {
				log.Printf("[%v] %v", s.Name(), err)
				return
			}
			go s.Listen()
		case SIGNAL_RESTART:
			s.Close()
			time.Sleep(2 * time.Second)
			if err := s.Connect(); err != nil {
				log.Printf("[%v] %v", s.Name(), err)
				return
			}
			go s.Listen()
		default:
		}
	}
}

// ticker is the normalized struct for use in topic 'tickers'
type Ticker struct {
	Exchange  string  `json:"x"`
	Pair      string  `json:"p"`
	Bid       float64 `json:"b"`
	BidSize   float64 `json:"bs"`
	Ask       float64 `json:"a"`
	AskSize   float64 `json:"as"`
	EventTime int64   `json:"et"` // unix time ms
}

func (t *Ticker) Fmt() (*Ticker, bool) {
	t.Pair = strings.ReplaceAll(t.Pair, "/", "")
	t.Pair = strings.ReplaceAll(t.Pair, "-", "")
	t.Pair = strings.ReplaceAll(t.Pair, ":", "")
	if i := strings.Index(t.Pair, "USD"); i > -1 {
		t.Pair = t.Pair[:i] + ":" + t.Pair[i:]
	}
	return t, t.Valid()
}

func (t *Ticker) Valid() bool {
	return len(t.Pair) > 0 && len(t.Exchange) > 0 && t.Bid > 0 && t.Ask > 0
}

func (t *Ticker) Key() string {
	return t.Exchange + ":" + t.Pair
}

func (t *Ticker) Bytes() []byte {
	b, _ := json.Marshal(t)
	return b
}
