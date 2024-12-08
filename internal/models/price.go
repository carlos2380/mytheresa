package models

type Price struct {
	Original int     `json:"original"`
	Final    int     `json:"value"`
	Discount *string `json:"discount_percentage"`
	Currency string  `json:"currency"`
}
