package product_application

import (
	"mytheresa/internal/models"
	"mytheresa/internal/storage"
	"strconv"
)

type ProductApplication struct {
	StgProduct  storage.ProductStoreage
	StgDiscount storage.DiscountStoreage
}

func (apProduct *ProductApplication) GetProducts(category string, priceLessThan *int, cursor string) ([]models.ProductPrice, string, error) {
	products, nextCursor, err := apProduct.StgProduct.GetProducts(category, priceLessThan, cursor)
	if err != nil {
		return nil, "", err
	}
	productsPrice, err := apProduct.Apply_discount(products)
	if err != nil {
		return nil, "", err
	}
	return productsPrice, nextCursor, nil
}

func (apProduct *ProductApplication) Apply_discount(products []models.Product) ([]models.ProductPrice, error) {
	var result []models.ProductPrice

	for _, product := range products {

		percent, err := apProduct.StgDiscount.GetDiscount(product.Category)
		if err != nil {
			return nil, err
		}

		if product.SkuDiscount != nil {
			if *product.SkuDiscount > percent {
				percent = *product.SkuDiscount
			}
		}

		var discountPercentage *string
		finalPrice := product.Price

		if percent > 0 {
			finalPrice = product.Price * (100 - percent) / 100
			discountStr := strconv.Itoa(percent) + "%"
			discountPercentage = &discountStr

		}

		productPrice := models.ProductPrice{
			SKU:      product.SKU,
			Name:     product.Name,
			Category: product.Category,
			Price: models.Price{
				Original: product.Price,
				Final:    finalPrice,
				Discount: discountPercentage,
				Currency: "EUR",
			},
		}

		result = append(result, productPrice)
	}

	return result, nil
}
