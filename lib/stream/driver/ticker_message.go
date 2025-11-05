package driver

import (
	"encoding/json"
	"fmt"
	"strings"
)

type TickerMessage struct {
	Exchange  string  `json:"x"`
	Pair      string  `json:"p"`
	Bid       float64 `json:"b"`
	BidSize   float64 `json:"bs"`
	Ask       float64 `json:"a"`
	AskSize   float64 `json:"as"`
	EventTime int64   `json:"et"` // unix time ms
}

func (t *TickerMessage) Fmt() {
	t.Pair = strings.ReplaceAll(t.Pair, "/", "")
	t.Pair = strings.ReplaceAll(t.Pair, "-", "")
	t.Pair = strings.ReplaceAll(t.Pair, ":", "")
	if i := strings.Index(t.Pair, "USD"); i > -1 {
		t.Pair = t.Pair[:i] + ":" + t.Pair[i:]
	}
}

// implements kafka.KMessage
func (t *TickerMessage) IsValid() bool {
	return t != nil && len(t.Pair) > 0 && len(t.Exchange) > 0 && t.Bid > 0 && t.Ask > 0
}

// implements kafka.KMessage
func (t *TickerMessage) Key() string {
	return fmt.Sprintf("%v:%v", t.Exchange, t.Pair)
}

// implements kafka.KMessage
func (t *TickerMessage) Bytes() []byte {
	b, _ := json.Marshal(t)
	return b
}
