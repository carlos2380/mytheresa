package storage

import "mytheresa/internal/models"

type ProductStoreage interface {
	GetProducts(category string, priceLessThan *int, cursor string) ([]models.Product, string, error)
}
