package v1

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/PavelParvadov/price_analyzer/price-producer/internal/domain/dto"
	"github.com/PavelParvadov/price_analyzer/price-producer/internal/domain/models"
	"github.com/PavelParvadov/price_analyzer/price-producer/internal/usecase"
	"go.uber.org/zap"
)

type Handler struct {
	publisher usecase.PricePublisher
	log       *zap.Logger
}

func NewHandler(p usecase.PricePublisher) *Handler {
	return &Handler{publisher: p}
}

func (h *Handler) ProducePrice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		if h.log != nil {
			h.log.Warn("bad method", zap.String("method", r.Method))
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	addPriceRequest := dto.AddPriceRequest{}

	if err := json.NewDecoder(r.Body).Decode(&addPriceRequest); err != nil {
		if h.log != nil {
			h.log.Warn("bad request body", zap.Error(err))
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validateBody(&addPriceRequest); err == false {
		if h.log != nil {
			h.log.Warn("validation failed", zap.String("symbol", addPriceRequest.Symbol), zap.Float64("value", addPriceRequest.Value))
		}
		http.Error(w, "Недостаточно аргументов", http.StatusBadRequest)
		return
	}
	symbol := strings.ToUpper(strings.TrimSpace(addPriceRequest.Symbol))
	price := models.Price{
		Symbol:    symbol,
		Value:     float64(addPriceRequest.Value),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	if h.log != nil {
		h.log.Info("incoming manual price", zap.String("symbol", price.Symbol), zap.Float64("value", price.Value))
	}
	err := h.publisher.Publish(r.Context(), price)
	if err != nil {
		if h.log != nil {
			h.log.Error("publish error", zap.Error(err))
		}
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if h.log != nil {
		h.log.Info("manual price published", zap.String("symbol", price.Symbol), zap.Float64("value", price.Value))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "price published successfully",
	})

}

func validateBody(d *dto.AddPriceRequest) bool {
	if d.Symbol == "" {
		return false
	}
	if d.Value < 0 {
		return false
	}
	return true
}
