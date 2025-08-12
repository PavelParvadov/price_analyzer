package grpc

import (
	"context"
	"time"

	"github.com/PavelParvadov/price_analyzer/price-cache/internal/usecase"
	pricepb "github.com/PavelParvadov/price_analyzer/price-cache/protoc/gen/go/price-cache"
)

type PriceGRPCServer struct {
	pricepb.UnimplementedPriceServiceServer
	uc *usecase.PriceUC
}

func NewPriceGRPCServer(uc *usecase.PriceUC) *PriceGRPCServer {
	return &PriceGRPCServer{uc: uc}
}

func (s *PriceGRPCServer) GetLatestPrice(ctx context.Context, req *pricepb.GetLatestPriceRequest) (*pricepb.GetLatestPriceResponse, error) {
	p, exists, err := s.uc.GetLatest(ctx, req.GetSymbol())
	if err != nil {
		return nil, err
	}

	var ts string
	if !p.Timestamp.IsZero() {
		ts = p.Timestamp.UTC().Format(time.RFC3339)
	}

	return &pricepb.GetLatestPriceResponse{
		Exists: exists,
		Price: &pricepb.Price{
			Symbol:    p.Symbol,
			Value:     p.Value,
			Timestamp: ts,
		},
	}, nil
}
