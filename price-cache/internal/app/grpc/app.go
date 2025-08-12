package grpcapp

import (
	"fmt"
	"net"

	delivery "github.com/PavelParvadov/price_analyzer/price-cache/internal/delivery/grpc"
	"github.com/PavelParvadov/price_analyzer/price-cache/internal/usecase"
	pricepb "github.com/PavelParvadov/price_analyzer/price-cache/protoc/gen/go/price-cache"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCApp struct {
	log  *zap.Logger
	grpc *grpc.Server
	port int
}

func NewGRPCApp(log *zap.Logger, port int, uc *usecase.PriceUC) *GRPCApp {
	srv := grpc.NewServer()
	priceServer := delivery.NewPriceGRPCServer(uc)
	pricepb.RegisterPriceServiceServer(srv, priceServer)
	return &GRPCApp{log: log, grpc: srv, port: port}
}

func (a *GRPCApp) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return err
	}
	a.log.Info("grpc server listening", zap.Int("port", a.port))
	return a.grpc.Serve(lis)
}

func (a *GRPCApp) Stop() {
	a.log.Info("grpc server stopping")
	a.grpc.GracefulStop()
}
