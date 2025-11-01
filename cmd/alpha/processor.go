package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/jimtang2/bc-go/lib/db"
	"github.com/jimtang2/bc-go/lib/stream"
)

type Processor struct {
	alphaChan chan *Alpha
	pairs     map[string]*PairMarket
	exchanges map[string]db.Exchange
}

func NewProcessor() *Processor {
	p := &Processor{
		alphaChan: make(chan *Alpha),
		pairs:     map[string]*PairMarket{},
		exchanges: map[string]db.Exchange{},
	}
	for _, e := range db.Exchanges() {
		p.exchanges[e.ID] = e
	}
	return p
}

func (p *Processor) Process(m *sarama.ConsumerMessage) {
	var ticker stream.Ticker
	if err := json.Unmarshal(m.Value, &ticker); err != nil {
		log.Println(err)
		return
	}
	if _, ok := p.pairs[ticker.Pair]; !ok {
		p.pairs[ticker.Pair] = &PairMarket{
			mu:        sync.Mutex{},
			name:      ticker.Pair,
			exchanges: p.exchanges,
		}
	}
	pair := p.pairs[ticker.Pair]
	if pair.lowestAsk == nil || ticker.Ask < pair.lowestAsk.Ask {
		pair.setLowestAsk(&ticker)
	}
	if pair.highestBid == nil || ticker.Bid > pair.highestBid.Bid {
		pair.setHighestBid(&ticker)
	}
	// find alpha
	if a := pair.spreadCalc(); a != nil {
		p.alphaChan <- a
	}
}

type PairMarket struct {
	mu         sync.Mutex
	exchanges  map[string]db.Exchange
	name       string
	lowestAsk  *stream.Ticker
	highestBid *stream.Ticker
}

func (p *PairMarket) spreadCalc() *Alpha {
	p.mu.Lock()
	if p.lowestAsk == nil || p.highestBid == nil || p.lowestAsk.Exchange == p.highestBid.Exchange {
		p.mu.Unlock()
		return nil
	}
	var (
		lo  = *p.lowestAsk
		hi  = *p.highestBid
		vol = 0.0
	)
	p.mu.Unlock()
	if hi.BidSize < lo.AskSize {
		vol = hi.BidSize
	} else {
		vol = lo.AskSize
	}
	// filter negative bid - ask prices
	if hi.Bid-lo.Ask <= 0 {
		return nil
	}
	return (&Alpha{
		AskExchange: lo.Exchange,
		AskPrice:    lo.Ask,
		AskSize:     lo.AskSize,
		AskFee:      p.exchanges[lo.Exchange].TakerFee,
		AskTime:     lo.EventTime,
		BidExchange: hi.Exchange,
		BidPrice:    hi.Bid,
		BidSize:     hi.BidSize,
		BidFee:      p.exchanges[hi.Exchange].TakerFee,
		BidTime:     hi.EventTime,
		Pair:        p.name,
		Spread:      hi.Bid - lo.Ask,
		Volume:      vol,
		Timestamp:   time.Now().UnixMilli(),
	}).profitCalc()
}

func (m *PairMarket) setLowestAsk(ticker *stream.Ticker) {
	m.mu.Lock()
	m.lowestAsk = ticker
	m.mu.Unlock()
	go m.void(ticker, 1*time.Second)
}

func (m *PairMarket) setHighestBid(ticker *stream.Ticker) {
	m.mu.Lock()
	m.highestBid = ticker
	m.mu.Unlock()
	go m.void(ticker, 1*time.Second)
}

func (m *PairMarket) void(ticker *stream.Ticker, d time.Duration) {
	time.Sleep(d)
	if m.lowestAsk == ticker {
		m.mu.Lock()
		m.lowestAsk = nil
		m.mu.Unlock()
	}
	if m.highestBid == ticker {
		m.mu.Lock()
		m.highestBid = nil
		m.mu.Unlock()
	}
}

type Alpha struct {
	Pair        string  `json:"p"`
	Spread      float64 `json:"s"`
	Volume      float64 `json:"v"`
	AskExchange string  `json:"ax"`
	AskPrice    float64 `json:"ap"`
	AskSize     float64 `json:"as"`
	AskFee      float64 `json:"af"`
	AskTime     int64   `json:"at"`
	BidExchange string  `json:"bx"`
	BidPrice    float64 `json:"bp"`
	BidSize     float64 `json:"bs"`
	BidFee      float64 `json:"bf"`
	BidTime     int64   `json:"bt"`
	Profit      float64 `json:"pl"`
	Timestamp   int64   `json:"ts"`
}

func (a *Alpha) profitCalc() *Alpha {
	a.Profit = a.Volume * (a.BidPrice*(1-a.BidFee) - a.AskPrice*(1+a.AskFee))
	return a
}

func (a *Alpha) Bytes() []byte {
	b, _ := json.Marshal(a)
	return b
}

func (a *Alpha) Key() string {
	return fmt.Sprintf("%s %v-%v $%.2f",
		a.Pair,
		a.BidExchange[:3],
		a.AskExchange[:3],
		a.Profit,
	)
}
