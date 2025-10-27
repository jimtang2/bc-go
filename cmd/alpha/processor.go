package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/jimtang2/bc-go/lib/stream"
)

type Processor struct {
	markets map[string]*Market // markets by pair name
	out     chan ProcessorOutput
}

type ProcessorOutput interface {
	Bytes() []byte
	Key() string
}

func NewProcessor() *Processor {
	return &Processor{
		out:     make(chan ProcessorOutput),
		markets: map[string]*Market{},
	}
}

func (p *Processor) Process(m *sarama.ConsumerMessage) {
	var t stream.Ticker
	if err := json.Unmarshal(m.Value, &t); err != nil {
		log.Println(err)
	}
	pair := t.Pair
	// lazy create market for pair
	if _, ok := p.markets[pair]; !ok {
		p.markets[pair] = &Market{
			mu:   sync.Mutex{},
			Pair: pair,
		}
	}
	market := p.markets[pair]
	// compare ticker to market lowest ask
	if market.LowestAsk == nil || t.Ask < market.LowestAsk.Ask {
		market.SetLowestAsk(&t)
	}
	// compare ticker to market highest bid
	if market.HighestBid == nil || t.Bid < market.HighestBid.Bid {
		market.SetHighestBid(&t)
	}
	if a := market.Seek(); a != nil {
		p.out <- a
	}
}

type Market struct {
	mu         sync.Mutex
	Pair       string
	LowestAsk  *stream.Ticker
	HighestBid *stream.Ticker
}

func (m *Market) SetLowestAsk(t *stream.Ticker) {
	m.mu.Lock()
	m.LowestAsk = t
	m.mu.Unlock()
	go m.Void(t, 1*time.Second)
}

func (m *Market) SetHighestBid(t *stream.Ticker) {
	m.mu.Lock()
	m.HighestBid = t
	m.mu.Unlock()
	go m.Void(t, 1*time.Second)
}

func (m *Market) Void(t *stream.Ticker, d time.Duration) {
	time.Sleep(d)
	if m.LowestAsk == t {
		m.mu.Lock()
		m.LowestAsk = nil
		m.mu.Unlock()
	}
	if m.HighestBid == t {
		m.mu.Lock()
		m.HighestBid = nil
		m.mu.Unlock()
	}
}

func (m *Market) Seek() *Alpha {
	if m.LowestAsk == nil || m.HighestBid == nil || m.LowestAsk.Exchange == m.HighestBid.Exchange {
		return nil
	}
	lo, hi := m.LowestAsk, m.HighestBid
	spread := hi.Bid - lo.Ask
	spreadPct := spread / lo.Ask * 100
	spreadSize := 0.0
	if hi.BidSize < lo.AskSize {
		spreadSize = hi.BidSize
	} else {
		spreadSize = lo.AskSize
	}
	spreadVal := spreadSize * spread
	if spread <= 0 {
		return nil
	}
	log.Printf("%10s  %v-%v  +%.2f%%  +%.2f$  %.2f  %.2f  [%.2f  %.2f  %.2f  %.2f]",
		m.Pair,
		hi.Exchange[:3],
		lo.Exchange[:3],
		spreadPct,
		spreadVal,
		spread,
		spreadSize,
		hi.Bid,
		lo.Ask,
		hi.BidSize,
		lo.AskSize,
	)
	return &Alpha{
		Ask:         *m.LowestAsk,
		Bid:         *m.HighestBid,
		Pair:        m.Pair,
		Spread:      spread,
		SpreadRatio: spread / lo.Ask,
		SpreadSize:  spreadSize,
	}
}

type Alpha struct {
	Ask         stream.Ticker `json:"ask"`
	Bid         stream.Ticker `json:"bid"`
	Pair        string        `json:"pair"`
	Spread      float64       `json:"spread"`
	SpreadRatio float64       `json:"spread_ratio"`
	SpreadSize  float64       `json:"spread_size"`
}

func (a *Alpha) Bytes() []byte {
	b, _ := json.Marshal(a)
	return b
}

func (a *Alpha) Key() string {
	return fmt.Sprintf("%s %v-%v $%.2f",
		a.Bid.Pair,
		a.Bid.Exchange[:3],
		a.Ask.Exchange[:3],
		a.Spread*a.SpreadSize,
	)
}
