package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	SKU         string             `bson:"sku" json:"sku"`
	Name        string             `bson:"name" json:"name"`
	Category    string             `bson:"category" json:"category"`
	Price       int                `bson:"price" json:"price"`
	SkuDiscount *int               `bson:"sku_discount,omitempty" json:"sku_discount,omitempty"`
}

type ProductPrice struct {
	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    Price  `json:"price"`
}
