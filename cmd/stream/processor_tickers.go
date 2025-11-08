package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/jimtang2/bc-go/lib/db"
	"github.com/jimtang2/bc-go/lib/kafka"
	"github.com/jimtang2/bc-go/pkg/pb/v1"
	"google.golang.org/protobuf/proto"
)

type TickersProcessor struct {
	highestBids sync.Map // string → *pb.Ticker
	lowestAsks  sync.Map // string → *pb.Ticker
	fees        map[string]float64
}

func NewTickersProcessor() *TickersProcessor {
	proc := &TickersProcessor{
		fees: db.ExchangeFees(),
	}
	return proc
}

func (proc *TickersProcessor) Process(m *sarama.ConsumerMessage) kafka.KMessage {
	var ticker pb.Ticker
	if err := proto.Unmarshal(m.Value, &ticker); err != nil {
		log.Println(err)
		return nil
	}
	if lo, loaded := proc.lowestAsks.LoadOrStore(ticker.Pair, &ticker); !loaded {
		// new ticker stored
	} else if ticker.Ask <= lo.(*pb.Ticker).Ask {
		// new ticker compare with loaded ticker
		proc.lowestAsks.Store(ticker.Pair, &ticker)
	}
	if hi, loaded := proc.highestBids.LoadOrStore(ticker.Pair, &ticker); !loaded {
		// new ticker stored
	} else if ticker.Bid >= hi.(*pb.Ticker).Bid {
		// new ticker compare with loaded ticker
		proc.highestBids.Store(ticker.Pair, &ticker)
	}
	go proc.expireTicker(&ticker)
	if match := proc.findMatch(ticker.Pair); match == nil {
		return nil
	} else {
		return &Match{match}
	}
}

func (proc *TickersProcessor) findMatch(pair string) *pb.Match {
	loI, loOk := proc.lowestAsks.Load(pair)
	hiI, hiOk := proc.highestBids.Load(pair)
	if !loOk || !hiOk {
		return nil
	}
	lo := loI.(*pb.Ticker)
	hi := hiI.(*pb.Ticker)
	if lo == hi {
		return nil
	}
	m := &pb.Match{
		Pair:        pair,
		AskExchange: lo.Exchange,
		AskPrice:    lo.Ask,
		AskSize:     lo.AskSize,
		AskTime:     lo.EventTime,
		AskFeeRate:  proc.fees[lo.Exchange],
		BidExchange: hi.Exchange,
		BidPrice:    hi.Bid,
		BidSize:     hi.BidSize,
		BidTime:     hi.EventTime,
		BidFeeRate:  proc.fees[hi.Exchange],
		Timestamp:   time.Now().UnixMilli(),
	}
	c := &pb.Calculations{}
	c.Volume = func(a, b float64) float64 {
		if a <= b {
			return a
		} else {
			return b
		}
	}(hi.BidSize, lo.BidSize)
	c.PriceDiff = hi.Bid - lo.Ask
	c.PriceAvg = (hi.Bid + lo.Ask) / 2
	c.Spread = c.PriceDiff * c.Volume
	c.SpreadPct = c.PriceDiff / c.PriceAvg
	c.BidFee = m.BidFeeRate * hi.Bid * c.Volume
	c.AskFee = m.AskFeeRate * lo.Ask * c.Volume
	c.ProfitLoss = c.Spread - c.BidFee - c.AskFee
	m.Calculations = c
	return m
}

func (proc *TickersProcessor) expireTicker(t *pb.Ticker) {
	time.Sleep(1 * time.Second)
	if lo, loaded := proc.lowestAsks.Load(t.Pair); loaded && lo == t {
		proc.lowestAsks.Delete(t.Pair)
	}
	if hi, loaded := proc.highestBids.Load(t.Pair); loaded && hi == t {
		proc.highestBids.Delete(t.Pair)
	}
}

type Match struct {
	*pb.Match
}

func (m *Match) IsValid() bool {
	return m.Match != nil
}

func (m *Match) Key() string {
	return fmt.Sprintf("%s %v-%v %.2f$",
		m.Match.Pair,
		m.Match.AskExchange[:3],
		m.Match.BidExchange[:3],
		m.Match.Calculations.ProfitLoss,
	)
}

func (m *Match) Bytes() []byte {
	b, _ := proto.Marshal(m.Match)
	return b
}
