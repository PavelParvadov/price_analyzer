package repository

import (
	"context"

	"github.com/PavelParvadov/price_analyzer/price-cache/internal/domain/models"
)

type PriceRepository interface {
	SaveLatest(ctx context.Context, price models.Price) error
	GetLatest(ctx context.Context, symbol string) (price models.Price, exists bool, err error)
}
