package app

import (
	"context"

	httpapp "github.com/PavelParvadov/price_analyzer/price-producer/internal/app/http"
	"github.com/PavelParvadov/price_analyzer/price-producer/internal/config"
	"github.com/PavelParvadov/price_analyzer/price-producer/internal/infrastructure/kafka"
	"github.com/PavelParvadov/price_analyzer/price-producer/internal/usecase"
	"go.uber.org/zap"
)

type App struct {
	HttpApp  *httpapp.HttpApp
	Producer *usecase.Producer
	Kafka    *kafka.Producer

	cancel context.CancelFunc
	log    *zap.Logger
}

func NewApp(log *zap.Logger, cfg config.Config) *App {

	kafkaProducer := kafka.NewProducer(cfg)
	uc := usecase.NewProducer(kafkaProducer, cfg)
	http := httpapp.NewHttpApp(cfg, kafkaProducer, log)

	return &App{
		HttpApp:  http,
		Producer: uc,
		Kafka:    kafkaProducer,
		log:      log,
	}
}

func (a *App) StartGenerator() {
	if a.cancel != nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel

	go func() {
		_ = a.Producer.Start(ctx)
	}()
}

func (a *App) Stop(ctx context.Context) {
	if a.cancel != nil {
		a.cancel()
	}
	if a.HttpApp != nil {
		a.HttpApp.Stop(ctx)
	}
	if a.Kafka != nil {
		_ = a.Kafka.Close()
	}
}
