package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/PavelParvadov/price_analyzer/price-producer/internal/config"
	"github.com/PavelParvadov/price_analyzer/price-producer/internal/domain/models"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	Producer *kafka.Writer
	log      *zap.Logger
}

func NewProducer(cfg config.Config, log *zap.Logger) *Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      cfg.KafkaConfig.Addresses,
		Topic:        cfg.KafkaConfig.Topic,
		Balancer:     &kafka.Hash{},
		RequiredAcks: 1,
		BatchTimeout: 50 * time.Millisecond,
		BatchSize:    100,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})
	return &Producer{
		Producer: writer,
		log:      log,
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
		if p.log != nil {
			p.log.Error("kafka publish failed", zap.Error(err), zap.String("symbol", pr.Symbol), zap.Float64("value", pr.Value))
		}
		return fmt.Errorf("error writing to kafka: %w", err)
	}
	if p.log != nil {
		p.log.Info("kafka message sent", zap.String("topic", p.Producer.Stats().Topic), zap.String("symbol", pr.Symbol), zap.Float64("value", pr.Value))
	}

	return nil

}

func (p *Producer) Close() error {
	return p.Producer.Close()
}
