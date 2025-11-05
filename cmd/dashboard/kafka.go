package main

import (
	"bytes"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/jimtang2/bc-go/lib/kafka"
)

// SpreadsProcessor just sends topic messages to the websocket http handler channel to be broadcasted to all active client websocket connections
type SpreadsProcessor struct {
	channel chan []byte
}

func NewSocketProxyProcessor(channel chan []byte) *SpreadsProcessor {
	return &SpreadsProcessor{
		channel: channel,
	}
}

// SpreadsProcessor Process to insert offset id to payload for frontend table row key
func (p *SpreadsProcessor) Process(m *sarama.ConsumerMessage) kafka.KMessage {
	b := bytes.Replace(m.Value, []byte("{"), []byte(fmt.Sprintf(`{"i":%v,`, m.Offset)), 1)
	p.channel <- b
	return nil
}
