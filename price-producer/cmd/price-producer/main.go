package main

import (
	"context"
	"github.com/PavelParvadov/price_analyzer/price-producer/internal/app"
	"github.com/PavelParvadov/price_analyzer/price-producer/internal/config"
	"github.com/PavelParvadov/price_analyzer/price-producer/pkg/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log := logger.NewLogger()
	log.Info("test", zap.Any("msg", "msg"))
	cfg := config.GetInstance()
	log.Info("config", zap.Any("cfg", cfg))
	application := app.NewApp(log, *cfg)
	application.StartGenerator()
	go application.HttpApp.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Info("stopping application")
	application.Stop(context.Background())
}
