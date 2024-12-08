package models

type DiscountCategory struct {
	Category string `json:"category" bson:"category"`
	Percent  int    `json:"percentage" bson:"percentage"`
}
