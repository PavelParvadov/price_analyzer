package v1

import (
	"net/http"
)

func (h *Handler) InitRoutes(mux *http.ServeMux) {
	mux.HandleFunc("v1/produce", h.ProducePrice)
}
