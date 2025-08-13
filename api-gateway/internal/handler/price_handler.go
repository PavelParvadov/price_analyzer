package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/PavelParvadov/price_analyzer/api-gateway/internal/service"
)

type PriceHandler struct {
	svc *service.PriceService
}

func NewPriceHandler(svc *service.PriceService) *PriceHandler {
	return &PriceHandler{svc: svc}
}

func (h *PriceHandler) GetLatestPrice(w http.ResponseWriter, r *http.Request) {
	symbol := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("symbol")))
	if symbol == "" {
		http.Error(w, "symbol is required", http.StatusBadRequest)
		return
	}
	resp, err := h.svc.GetLatestPrice(r.Context(), symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
