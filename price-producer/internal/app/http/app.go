package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/PavelParvadov/price_analyzer/price-producer/internal/config"
	v1 "github.com/PavelParvadov/price_analyzer/price-producer/internal/delivery/http/v1"
	"github.com/PavelParvadov/price_analyzer/price-producer/internal/usecase"
	"go.uber.org/zap"
)

type HttpApp struct {
	Server *http.Server
	Log    *zap.Logger
}

func NewHttpApp(cfg config.Config, publisher usecase.PricePublisher, log *zap.Logger) *HttpApp {
	mux := http.NewServeMux()
	h := v1.NewHandler(publisher)
	h.InitRoutes(mux)
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.PriceConfig.Port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second, //TODO вынести в конфиг
	}
	return &HttpApp{
		Server: srv,
		Log:    log,
	}
}

func (app *HttpApp) Run() error {
	app.Log.Info("http server listening", zap.String("addr", app.Server.Addr))
	if err := app.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (app *HttpApp) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *HttpApp) Stop(ctx context.Context) {
	app.Log.Info("http server stopping")
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_ = app.Server.Shutdown(shutdownCtx)
}
