package main

import (
	httpapp "github.com/PavelParvadov/price_analyzer/api-gateway/internal/app/http"
	"github.com/PavelParvadov/price_analyzer/api-gateway/internal/config"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log, _ := zap.NewDevelopment()
	cfg := config.GetInstance()

	app, err := httpapp.NewHTTPApp(log, *cfg)
	if err != nil {
		log.Fatal("failed to init http app", zap.Error(err))
	}
	go app.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	app.Stop()

}
