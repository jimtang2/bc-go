package main

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/jimtang2/bc-go/lib/db"
	"github.com/jimtang2/bc-go/lib/kafka"
	"github.com/jimtang2/bc-go/pkg/pb/v1"
	"google.golang.org/protobuf/proto"
)

type AlphaProcessor struct{}

func (p *AlphaProcessor) Process(m *sarama.ConsumerMessage) kafka.KMessage {
	v := pb.Match{}
	if err := proto.Unmarshal(m.Value, &v); err != nil {
		log.Println(err)
		return nil
	}
	log.Println(&v)
	v.Id = m.Offset
	if err := db.InsertMatch(&v); err != nil {
		log.Println(err)
	}
	return nil
}
