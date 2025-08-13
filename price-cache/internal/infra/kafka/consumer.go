package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"go.uber.org/zap"

	cachecfg "github.com/PavelParvadov/price_analyzer/price-cache/internal/config"
	"github.com/PavelParvadov/price_analyzer/price-cache/internal/domain/models"
	"github.com/PavelParvadov/price_analyzer/price-cache/internal/usecase"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	log    *zap.Logger
}

func NewConsumer(cfg cachecfg.Config, log *zap.Logger) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:               cfg.Kafka.Addresses,
		GroupID:               cfg.Kafka.GroupID,
		Topic:                 cfg.Kafka.Topic,
		MinBytes:              1,
		MaxBytes:              10e6,
		StartOffset:           kafka.FirstOffset,
		ReadLagInterval:       0,
		WatchPartitionChanges: true,
		HeartbeatInterval:     3 * time.Second,
		SessionTimeout:        30 * time.Second,
		RebalanceTimeout:      30 * time.Second,
		CommitInterval:        time.Second,
	})
	return &Consumer{reader: r, log: log}
}

func (c *Consumer) Start(ctx context.Context, uc *usecase.PriceUC) error {
	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return err
			}
			time.Sleep(500 * time.Millisecond)
			continue
		}

		var p models.Price
		if err := json.Unmarshal(m.Value, &p); err == nil {
			c.log.Info("kafka message received", zap.String("symbol", p.Symbol), zap.Float64("value", p.Value))
			err = uc.SaveLatest(ctx, p.Symbol, p.Value, p.Timestamp)
			if err != nil {
				c.log.Error("failed to save price", zap.Error(err))
			}
		} else {
			c.log.Warn("failed to unmarshal message", zap.Error(err))
		}

		err = c.reader.CommitMessages(ctx, m)
		if err != nil {
			c.log.Error("failed to commit message", zap.Error(err))
		} else {
			c.log.Info("kafka message committed", zap.Int("partition", m.Partition), zap.Int64("offset", m.Offset))
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
