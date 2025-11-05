package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/jimtang2/bc-go/lib/db"
	"github.com/jimtang2/bc-go/lib/kafka"
	"github.com/jimtang2/bc-go/lib/stream/driver"
)

type TickersProcessor struct {
	pairs     map[string]*PairMarket
	exchanges map[string]db.Exchange
}

func NewTickersProcessor() *TickersProcessor {
	p := &TickersProcessor{
		pairs:     map[string]*PairMarket{},
		exchanges: map[string]db.Exchange{},
	}
	for _, e := range db.Exchanges() {
		p.exchanges[e.ID] = e
	}
	return p
}

func (p *TickersProcessor) Process(m *sarama.ConsumerMessage) kafka.KMessage {
	var ticker driver.TickerMessage
	if err := json.Unmarshal(m.Value, &ticker); err != nil {
		log.Println(err)
		return nil
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
	return pair.spreadCalc()
}

type PairMarket struct {
	mu         sync.Mutex
	exchanges  map[string]db.Exchange
	name       string
	lowestAsk  *driver.TickerMessage // lowest price A will sell
	highestBid *driver.TickerMessage // highest price A will buy
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

func (m *PairMarket) setLowestAsk(ticker *driver.TickerMessage) {
	m.mu.Lock()
	m.lowestAsk = ticker
	m.mu.Unlock()
	go m.setExpire(ticker, 1*time.Second)
}

func (m *PairMarket) setHighestBid(ticker *driver.TickerMessage) {
	m.mu.Lock()
	m.highestBid = ticker
	m.mu.Unlock()
	go m.setExpire(ticker, 1*time.Second)
}

func (m *PairMarket) setExpire(ticker *driver.TickerMessage, d time.Duration) {
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
	var (
		v1 = a.Volume * a.BidPrice
		v2 = a.Volume * a.AskPrice
		v3 = (a.Volume * a.BidPrice) * a.BidFee / 100
		v4 = (a.Volume * a.AskPrice) * a.AskFee / 100
	)
	a.Profit = v1 - v2 - v3 - v4
	return a
}

func (a *Alpha) Bytes() []byte {
	b, _ := json.Marshal(a)
	return b
}

func (a *Alpha) Key() string {
	return fmt.Sprintf("%s %v-%v %.2f$",
		a.Pair,
		a.BidExchange[:3],
		a.AskExchange[:3],
		a.Profit,
	)
}

func (a *Alpha) IsValid() bool {
	return a != nil
}
