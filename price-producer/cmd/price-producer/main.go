package main

import (
	"go.uber.org/zap"
	"price_analyzer/price-producer/pkg/logger"
)

func main() {
	log := logger.NewLogger()
	log.Info("test", zap.Any("any", "pavel"))
}
