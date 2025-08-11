package usecase

import (
	"context"
	"github.com/PavelParvadov/price_analyzer/price-producer/internal/config"
	"math/rand"
	"strings"
	"time"

	"github.com/PavelParvadov/price_analyzer/price-producer/internal/domain/models"
)

type PricePublisher interface {
	Publish(ctx context.Context, p models.Price) error
}

type Producer struct {
	Publisher PricePublisher
	Cfg       config.ProducerConfig
	Last      map[string]float64
}

func NewProducer(publisher PricePublisher, cfg config.ProducerConfig) *Producer {
	normalized := make([]string, 0, len(cfg.Tickers))
	for _, t := range cfg.Tickers {
		s := strings.ToUpper(strings.TrimSpace(t))
		if s != "" {
			normalized = append(normalized, s)
		}
	}
	cfg.Tickers = normalized

	return &Producer{
		Publisher: publisher,
		Cfg:       cfg,
		Last:      make(map[string]float64),
	}
}

func (p *Producer) Start(ctx context.Context) error {
	rand.Seed(time.Now().UnixNano())
	ticker := time.NewTicker(time.Duration(p.Cfg.IntervalMs))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			for _, symbol := range p.Cfg.Tickers {
				value := p.nextPrice(symbol)
				msg := models.Price{
					Symbol:    symbol,
					Value:     value,
					Timestamp: time.Now().UTC().Format(time.RFC3339),
				}
				_ = p.Publisher.Publish(ctx, msg)
			}
		}
	}
}

func (p *Producer) nextPrice(symbol string) float64 {
	prev := p.Last[symbol]
	if prev <= 0 {
		prev = p.Cfg.InitialPrice
	}
	vol := p.Cfg.VolatilityPercent / 100.0
	if vol < 0 {
		vol = 0
	}
	delta := (rand.Float64()*2 - 1) * vol
	next := prev * (1 + delta)
	p.Last[symbol] = next
	return next
}
