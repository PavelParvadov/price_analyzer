package usecase

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/PavelParvadov/price_analyzer/price-producer/internal/config"

	"github.com/PavelParvadov/price_analyzer/price-producer/internal/domain/models"
	"go.uber.org/zap"
)

type PricePublisher interface {
	Publish(ctx context.Context, p models.Price) error
}

type Producer struct {
	publisher PricePublisher
	cfg       config.Config
	last      map[string]float64
	log       *zap.Logger
}

func NewProducer(publisher PricePublisher, cfg config.Config) *Producer {
	normalized := make([]string, 0, len(cfg.Producer.Tickers))
	for _, t := range cfg.Producer.Tickers {
		s := strings.ToUpper(strings.TrimSpace(t))
		if s != "" {
			normalized = append(normalized, s)
		}
	}
	cfg.Producer.Tickers = normalized

	return &Producer{
		publisher: publisher,
		cfg:       cfg,
		last:      make(map[string]float64),
		log:       zap.NewNop(),
	}
}

func (p *Producer) Start(ctx context.Context) error {
	rand.Seed(time.Now().UnixNano())
	ticker := time.NewTicker(time.Duration(p.cfg.Producer.IntervalMs) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			p.Ticketing(ctx)
		}
	}

}

func (p *Producer) Ticketing(ctx context.Context) {
	for _, symbol := range p.cfg.Producer.Tickers {
		value := p.nextPrice(symbol)
		msg := models.Price{
			Symbol:    symbol,
			Value:     value,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		pubCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		if err := p.publisher.Publish(pubCtx, msg); err != nil {
			p.log.Error("publish failed", zap.String("symbol", msg.Symbol), zap.Float64("value", msg.Value), zap.Error(err))
		}
		cancel()
	}
}

func (p *Producer) nextPrice(symbol string) float64 {
	prev := p.last[symbol]
	if prev <= 0 {
		prev = p.cfg.Producer.InitialPrice
	}
	vol := p.cfg.Producer.VolatilityPercent / 100.0
	if vol < 0 {
		vol = 0
	}
	delta := (rand.Float64()*2 - 1) * vol
	next := prev * (1 + delta)
	p.last[symbol] = next
	return next
}
