package models

type Product struct {
	Id          string `json:"id"`
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Price       int    `json:"price"`
	SkuDiscount *int   `json:"sku_discount,omitempty"`
}

type ProductPrice struct {
	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    Price  `json:"price"`
}
