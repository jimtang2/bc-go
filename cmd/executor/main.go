// webhook$$bc-go;cmd/executor/main.go;grok$$
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

type Executor struct {
	consumer sarama.ConsumerGroup
	config   *viper.Viper
}

func NewExecutor() *Executor {
	v := viper.New()
	v.SetConfigFile("/config/config.yml")
	v.ReadInConfig()
	return &Executor{config: v}
}

func (e *Executor) Start(ctx context.Context) error {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	consumer, err := sarama.NewConsumerGroup([]string{e.config.GetString("kafka_brokers")}, "executor-group", config)
	if err != nil {
		return err
	}
	e.consumer = consumer
	
	handler := &executorHandler{executor: e}
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			err := consumer.Consume(ctx, []string{"trades"}, handler)
			if err != nil {
				log.Printf("Consume error: %v", err)
			}
		}
	}
}

type executorHandler struct {
	executor *Executor
}

func (h *executorHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *executorHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *executorHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		session.MarkMessage(msg, "")
	}
	return nil
}

func main() {
	e := NewExecutor()
	ctx, cancel := context.WithCancel(context.Background())
	
	go e.Start(ctx)
	
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	cancel()
	e.consumer.Close()
}