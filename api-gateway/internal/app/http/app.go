package httpapp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PavelParvadov/price_analyzer/api-gateway/internal/config"
	"github.com/PavelParvadov/price_analyzer/api-gateway/internal/handler"
	"github.com/PavelParvadov/price_analyzer/api-gateway/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type HTTPApp struct {
	srv *http.Server
	log *zap.Logger
}

func NewHTTPApp(log *zap.Logger, cfg config.Config) (*HTTPApp, error) {
	addr := fmt.Sprintf("%s:%d", cfg.PriceCache.Host, cfg.PriceCache.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	svc := service.NewPriceService(conn)
	h := handler.NewPriceHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/price", h.GetLatestPrice)

	srv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.HTTP.Port), Handler: mux}
	return &HTTPApp{srv: srv, log: log}, nil
}

func (a *HTTPApp) Run() error {
	a.log.Info("api-gateway http listening", zap.String("addr", a.srv.Addr))
	return a.srv.ListenAndServe()
}

func (a *HTTPApp) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *HTTPApp) Stop() {
	a.log.Info("api-gateway http shutdown", zap.String("addr", a.srv.Addr))
	err := a.srv.Shutdown(context.Background())
	if err != nil {
		return
	}
}
