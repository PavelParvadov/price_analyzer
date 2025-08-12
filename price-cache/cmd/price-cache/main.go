package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/PavelParvadov/price_analyzer/price-cache/internal/app"
	"github.com/PavelParvadov/price_analyzer/price-cache/internal/config"
	"go.uber.org/zap"
)

func main() {
	log, _ := zap.NewDevelopment()
	cfg := config.GetInstance()

	application := app.NewApp(log, *cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stop
		cancel()
		application.Stop()
	}()

	if err := application.Run(ctx); err != nil {
		log.Error("application stopped", zap.Error(err))
	}
}
