package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	cachecfg "github.com/PavelParvadov/price_analyzer/price-cache/internal/config"
	"github.com/PavelParvadov/price_analyzer/price-cache/internal/domain/models"
	"github.com/PavelParvadov/price_analyzer/price-cache/internal/repository"
	"github.com/redis/go-redis/v9"
)

type PriceRedisRepo struct {
	rdb *redis.Client
}

func NewPriceRedisRepo(cfg cachecfg.Config) *PriceRedisRepo {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	return &PriceRedisRepo{rdb: client}
}

func (r *PriceRedisRepo) key(symbol string) string {
	s := strings.ToUpper(strings.TrimSpace(symbol))
	return fmt.Sprintf("price:%s", s)
}

func (r *PriceRedisRepo) SaveLatest(ctx context.Context, price models.Price) error {

	data, err := json.Marshal(price)
	if err != nil {
		return err
	}
	return r.rdb.Set(ctx, r.key(price.Symbol), data, 0).Err()
}

func (r *PriceRedisRepo) GetLatest(ctx context.Context, symbol string) (models.Price, bool, error) {
	res, err := r.rdb.Get(ctx, r.key(symbol)).Result()
	if errors.Is(err, redis.Nil) {
		return models.Price{}, false, nil
	}
	if err != nil {
		return models.Price{}, false, err
	}
	var p models.Price
	if err := json.Unmarshal([]byte(res), &p); err != nil {
		return models.Price{}, false, err
	}
	return p, true, nil
}

func (r *PriceRedisRepo) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return r.rdb.Ping(ctx).Err()
}

var _ repository.PriceRepository = (*PriceRedisRepo)(nil)
