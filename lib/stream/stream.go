package stream

import (
	"encoding/json"
	"strings"
)

type Message struct {
	Topic   string
	Key     string
	Headers map[string]string
	Payload []byte
}

type Stream interface {
	Start() chan Message
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

func (t *Ticker) Fmt() *Ticker {
	t.Pair = strings.ReplaceAll(t.Pair, "/", "")
	t.Pair = strings.ReplaceAll(t.Pair, "-", "")
	t.Pair = strings.ReplaceAll(t.Pair, ":", "")
	if i := strings.Index(t.Pair, "USD"); i > -1 {
		t.Pair = t.Pair[:i] + ":" + t.Pair[i:]
	}
	return t
}

func (t *Ticker) IsValid() bool {
	return len(t.Pair) > 0 && len(t.Exchange) > 0 && t.Bid > 0 && t.Ask > 0
}

func (t *Ticker) Key() string {
	return t.Exchange + ":" + t.Pair
}

func (t *Ticker) Bytes() []byte {
	b, _ := json.Marshal(t)
	return b
}
