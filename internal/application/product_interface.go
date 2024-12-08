package application

import "mytheresa/internal/models"

type ProductIntef interface {
	GetProducts(category string, priceLessThan *int, cursor string) ([]models.ProductPrice, string, error)
}
