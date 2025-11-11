package main

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/jimtang2/bc-go/lib/kafka"
	"github.com/jimtang2/bc-go/pkg/pb/v1"
	"google.golang.org/protobuf/proto"
)

// MatchesProcessor just sends topic messages to the websocket http handler channel to be broadcasted to all active client websocket connections
type MatchesProcessor struct {
	channel chan []byte
}

func NewMatchesProxyProcessor(channel chan []byte) *MatchesProcessor {
	return &MatchesProcessor{
		channel: channel,
	}
}

// MatchesProcessor Process to insert offset id to payload for frontend table row key
func (p *MatchesProcessor) Process(m *sarama.ConsumerMessage) kafka.KMessage {
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

type TickersProcessor struct {
	channel chan []byte
}

func NewTickersProxyProcessor(channel chan []byte) *TickersProcessor {
	return &TickersProcessor{
		channel: channel,
	}
}

func (p *TickersProcessor) Process(m *sarama.ConsumerMessage) kafka.KMessage {
	p.channel <- m.Value
	return nil
}
