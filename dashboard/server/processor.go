package main

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/jimtang2/bc-go/lib/db"
	"github.com/jimtang2/bc-go/lib/stream"
)

type Processor struct {
	exchanges map[string]db.Exchange
}

func NewProcessor() *Processor {
	p := Processor{
		exchanges: map[string]db.Exchange{},
	}
	for _, e := range db.Exchanges() {
		p.exchanges[e.ID] = e
	}
	return &p
}

func (p *Processor) process(m *sarama.ConsumerMessage) ([]byte, error) {
	v := MessageIncoming{}
	if err := json.Unmarshal(m.Value, &v); err != nil {
		return []byte{}, err
	}
	return json.Marshal(MessageOutgoing{
		ID:          m.Offset,
		Pair:        v.Pair,
		Spread:      v.Spread,
		SpreadRatio: v.SpreadRatio,
		SpreadSize:  v.SpreadSize,
		AskExchange: p.exchanges[v.Ask.Exchange].Name,
		BidExchange: p.exchanges[v.Bid.Exchange].Name,
		AskPrice:    v.Ask.Ask,
		BidPrice:    v.Bid.Bid,
		AskSize:     v.Ask.AskSize,
		BidSize:     v.Bid.BidSize,
		AskTime:     v.Ask.EventTime,
		BidTime:     v.Bid.EventTime,
		AskFee:      p.exchanges[v.Ask.Exchange].TakerFee,
		BidFee:      p.exchanges[v.Bid.Exchange].TakerFee,
		ExpPL:       p.calcProfit(v.Ask, v.Bid),
	})
}

func (p *Processor) calcProfit(lo, hi stream.Ticker) float64 {
	var vol float64
	if lo.AskSize > hi.BidSize {
		vol = hi.BidSize
	} else {
		vol = lo.AskSize
	}
	var (
		min    = lo.Ask
		max    = hi.Bid
		feeMax = p.exchanges[hi.Exchange].TakerFee
		feeMin = p.exchanges[lo.Exchange].TakerFee
		fee    = max*vol*feeMax + min*vol*feeMin
		x      = vol*(max-min) - fee
	)
	return x
}

type MessageIncoming struct {
	Ask         stream.Ticker `json:"ask"`
	Bid         stream.Ticker `json:"bid"`
	Pair        string        `json:"pair"`
	Spread      float64       `json:"spread"`
	SpreadRatio float64       `json:"spread_ratio"`
	SpreadSize  float64       `json:"spread_size"`
}

type MessageOutgoing struct {
	ID          int64   `json:"i"`
	Pair        string  `json:"p"`
	Spread      float64 `json:"s"`
	SpreadRatio float64 `json:"sr"`
	SpreadSize  float64 `json:"ss"`
	AskExchange string  `json:"ax"`
	BidExchange string  `json:"bx"`
	AskPrice    float64 `json:"ap"`
	BidPrice    float64 `json:"bp"`
	AskSize     float64 `json:"as"`
	BidSize     float64 `json:"bs"`
	AskTime     int64   `json:"at"`
	BidTime     int64   `json:"bt"`
	AskFee      float64 `json:"af"`
	BidFee      float64 `json:"bf"`
	ExpPL       float64 `json:"pl"`
}
