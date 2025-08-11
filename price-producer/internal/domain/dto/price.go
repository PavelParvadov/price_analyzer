package dto

type AddPriceRequest struct {
	Symbol string  `json:"symbol"`
	Value  float32 `json:"value"`
}
