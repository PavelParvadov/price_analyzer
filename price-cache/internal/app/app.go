package app

import (
	"context"

	grpcapp "github.com/PavelParvadov/price_analyzer/price-cache/internal/app/grpc"
	"github.com/PavelParvadov/price_analyzer/price-cache/internal/config"
	kconsumer "github.com/PavelParvadov/price_analyzer/price-cache/internal/infra/kafka"
	rredis "github.com/PavelParvadov/price_analyzer/price-cache/internal/infra/redis"
	"github.com/PavelParvadov/price_analyzer/price-cache/internal/usecase"
	"go.uber.org/zap"
)

type App struct {
	log      *zap.Logger
	cfg      config.Config
	consumer *kconsumer.Consumer
	uc       *usecase.PriceUC
	grpcSrv  *grpcapp.GRPCApp
}

func NewApp(log *zap.Logger, cfg config.Config) *App {
	repo := rredis.NewPriceRedisRepo(cfg)
	uc := usecase.NewPriceUC(repo)
	consumer := kconsumer.NewConsumer(cfg, log)

	grpcSrv := grpcapp.NewGRPCApp(log, cfg.GRPC.Port, uc)

	return &App{
		log:      log,
		cfg:      cfg,
		consumer: consumer,
		uc:       uc,
		grpcSrv:  grpcSrv,
	}
}

func (a *App) Run(ctx context.Context) error {
	a.log.Info("price-cache: starting consumer",
		zap.Strings("brokers", a.cfg.Kafka.Addresses),
		zap.String("topic", a.cfg.Kafka.Topic),
		zap.String("group", a.cfg.Kafka.GroupID),
	)
	go func() {
		_ = a.consumer.Start(ctx, a.uc)
	}()
	return a.grpcSrv.Run()
}

func (a *App) Stop() {
	_ = a.consumer.Close()
	if a.grpcSrv != nil {
		a.grpcSrv.Stop()
	}
}
