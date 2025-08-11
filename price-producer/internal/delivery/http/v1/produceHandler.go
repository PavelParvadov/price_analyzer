package v1

import (
	"encoding/json"
	"github.com/PavelParvadov/price_analyzer/price-producer/internal/domain/dto"
	"github.com/PavelParvadov/price_analyzer/price-producer/internal/domain/models"
	"github.com/PavelParvadov/price_analyzer/price-producer/internal/usecase"
	"net/http"
	"time"
)

type Handler struct {
	publisher usecase.PricePublisher
}

func NewHandler(p usecase.PricePublisher) *Handler {
	return &Handler{publisher: p}
}

func (h *Handler) ProducePrice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	addPriceRequest := dto.AddPriceRequest{}

	if err := json.NewDecoder(r.Body).Decode(&addPriceRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validateBody(&addPriceRequest); err == false {
		http.Error(w, "Недостаточно аргументов", http.StatusBadRequest)
		return
	}
	price := models.Price{
		Symbol:    addPriceRequest.Symbol,
		Value:     float64(addPriceRequest.Value),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	err := h.publisher.Publish(r.Context(), price)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
