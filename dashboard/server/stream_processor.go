package main

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/jimtang2/bc-go/lib/db"
)

type Processor interface {
	Process(m *sarama.ConsumerMessage) error
}

// SpreadsProcessor just sends topic messages to the websocket http handler channel to be broadcasted to all active client websocket connections
type SpreadsProcessor struct {
	channel chan []byte
}

func NewSpreadsProcessor(channel chan []byte) *SpreadsProcessor {
	return &SpreadsProcessor{
		channel: channel,
	}
}

func (p *SpreadsProcessor) Process(m *sarama.ConsumerMessage) error {
	p.channel <- m.Value
	return nil
}

// AlphaProcessor pulls all records from topic 'alpha' (from oldest) and stores them in sql so dashboard client can display and query them
type AlphaProcessor struct{}

func NewAlphaProcessor() *AlphaProcessor {
	return &AlphaProcessor{}
}

func (p *AlphaProcessor) Process(m *sarama.ConsumerMessage) error {
	v := db.Alpha{}
	if err := json.Unmarshal(m.Value, &v); err != nil {
		return err
	}
	return db.InsertAlpha(v)
}
