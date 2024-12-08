package models

type DiscountCategory struct {
	Category string `json:"category"`
	Percent  int    `json:"percentage"`
}
