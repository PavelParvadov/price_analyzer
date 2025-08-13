package service

import (
	"context"

	pricepb "github.com/PavelParvadov/price_analyzer/api-gateway/protoc/gen/go/api-gateway"
	"google.golang.org/grpc"
)

type PriceService struct {
	client pricepb.PriceServiceClient
}

func NewPriceService(conn *grpc.ClientConn) *PriceService {
	return &PriceService{client: pricepb.NewPriceServiceClient(conn)}
}

func (s *PriceService) GetLatestPrice(ctx context.Context, symbol string) (*pricepb.GetLatestPriceResponse, error) {
	return s.client.GetLatestPrice(ctx, &pricepb.GetLatestPriceRequest{Symbol: symbol})
}
