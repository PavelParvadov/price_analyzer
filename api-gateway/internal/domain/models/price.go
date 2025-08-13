package models

type Price struct {
	Symbol    string  `json:"symbol"`
	Value     float64 `json:"value"`
	Timestamp string  `json:"timestamp"`
}
