package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/PavelParvadov/price_analyzer/price-cache/internal/domain/models"
	"github.com/PavelParvadov/price_analyzer/price-cache/internal/repository"
)

type PriceUC struct {
	repo repository.PriceRepository
}

func NewPriceUC(repo repository.PriceRepository) *PriceUC {
	return &PriceUC{repo: repo}
}

func (u *PriceUC) SaveLatest(ctx context.Context, symbol string, value float64, ts time.Time) error {
	s := strings.ToUpper(strings.TrimSpace(symbol))
	if s == "" {
		return nil
	}
	p := models.Price{Symbol: s, Value: value, Timestamp: ts.UTC()}
	return u.repo.SaveLatest(ctx, p)
}

func (u *PriceUC) GetLatest(ctx context.Context, symbol string) (models.Price, bool, error) {
	s := strings.ToUpper(strings.TrimSpace(symbol))
	if s == "" {
		return models.Price{}, false, nil
	}
	return u.repo.GetLatest(ctx, s)
}
