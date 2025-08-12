package models

import "time"

type Price struct {
	Symbol    string
	Value     float64
	Timestamp time.Time
}
