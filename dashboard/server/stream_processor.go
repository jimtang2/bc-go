package main

import (
	"bytes"
	"encoding/json"
	"fmt"

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

// SpreadsProcessor Process to insert offset id to payload for frontend table row key
func (p *SpreadsProcessor) Process(m *sarama.ConsumerMessage) error {
	b := bytes.Replace(m.Value, []byte("{"), []byte(fmt.Sprintf(`{"i":%v,`, m.Offset)), 1)
	p.channel <- b
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
	v.ID = m.Offset
	return db.InsertAlpha(v)
}
