package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PavelParvadov/price_analyzer/price-producer/internal/config"
	"github.com/PavelParvadov/price_analyzer/price-producer/internal/domain/models"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Producer *kafka.Writer
}

func NewProducer(cfg config.Config) *Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Topic:    cfg.KafkaConfig.Topic,
		Brokers:  cfg.KafkaConfig.Addresses,
		Balancer: &kafka.Hash{},
	})
	return &Producer{
		Producer: writer,
	}
}

func (p *Producer) Publish(ctx context.Context, pr models.Price) error {
	data, err := json.Marshal(&pr)
	if err != nil {
		return err
	}
	err = p.Producer.WriteMessages(ctx, kafka.Message{
		Value: data,
		Key:   []byte(pr.Symbol),
	})
	if err != nil {
		return fmt.Errorf("error writing to kafka: %w", err)
	}

	return nil

}

func (p *Producer) Close() error {
	return p.Producer.Close()
}
