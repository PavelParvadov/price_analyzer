package dto

type AddPriceRequest struct {
	Symbol string  `json:"symbol"`
	Value  float64 `json:"value"`
}
