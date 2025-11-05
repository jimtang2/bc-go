package main

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/jimtang2/bc-go/lib/db"
	"github.com/jimtang2/bc-go/lib/kafka"
)

type AlphaProcessor struct {
	pairs     map[string]*PairMarket
	exchanges map[string]db.Exchange
}

func NewAlphaProcessor() *AlphaProcessor {
	p := &AlphaProcessor{
		pairs:     map[string]*PairMarket{},
		exchanges: map[string]db.Exchange{},
	}
	for _, e := range db.Exchanges() {
		p.exchanges[e.ID] = e
	}
	return p
}

func (p *AlphaProcessor) Process(m *sarama.ConsumerMessage) kafka.KMessage {
	v := db.Alpha{}
	if err := json.Unmarshal(m.Value, &v); err != nil {
		log.Println(err)
		return nil
	}
	v.ID = m.Offset
	if err := db.InsertAlpha(v); err != nil {
		log.Println(err)
	}
	return nil
}
