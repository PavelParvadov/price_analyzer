package main

import (
	"github.com/PavelParvadov/price_analyzer/price-producer/internal/config"
	"github.com/PavelParvadov/price_analyzer/price-producer/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	log := logger.NewLogger()
	log.Info("test", zap.Any("msg", "msg"))
	cfg := config.GetInstance()
	log.Info("config", zap.Any("cfg", cfg))
}
