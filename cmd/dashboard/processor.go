package main

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/jimtang2/bc-go/lib/kafka"
	"github.com/jimtang2/bc-go/pkg/pb/v1"
	"google.golang.org/protobuf/proto"
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
	v := pb.Match{}
	err := proto.Unmarshal(m.Value, &v)
	if err != nil {
		log.Println(err)
		return nil
	}
	v.Id = m.Offset
	b, err := proto.Marshal(&v)
	if err != nil {
		log.Println(err)
		return nil
	}
	p.channel <- b
	return nil
}
