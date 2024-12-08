package mock

import (
	"encoding/json"
	"fmt"
	"log"
	"mytheresa/internal/models"
	"sort"
	"strconv"
)

type MockProduct struct {
	Products   []*models.Product
	ProductMap map[string][]*models.Product
}

func (mProduct *MockProduct) LoadData() error {

	err := json.Unmarshal([]byte(rawProductsJSON), &mProduct.Products)
	if err != nil {
		log.Println("Error decoding JSON:", err)
	}

	sort.Slice(mProduct.Products, func(i, j int) bool {
		return mProduct.Products[i].Price < mProduct.Products[j].Price
	})

	mProduct.ProductMap = make(map[string][]*models.Product)
	for i := 0; i < len(mProduct.Products); i++ {
		mProduct.ProductMap[mProduct.Products[i].Category] = append(mProduct.ProductMap[mProduct.Products[i].Category], mProduct.Products[i])
	}

	// for _, product := range mProduct.Products {
	// 	fmt.Println(product)
	// }
	return nil
}

func (mProduct *MockProduct) GetProducts(category string, priceLessThan *int, cursor string) ([]models.Product, string, error) {
	var productsReturn []*models.Product

	if category == "" {
		productsReturn = mProduct.Products
	} else {
		productsReturn = mProduct.ProductMap[category]
	}

	startIndex := 0
	if cursor != "" {
		var err error
		startIndex, err = strconv.Atoi(cursor)
		if err != nil || startIndex < 0 || startIndex >= len(productsReturn) {
			return nil, "", fmt.Errorf("cursor inválido")
		}
	}

	filteredProducts := []models.Product{}

	for i := startIndex; i < len(productsReturn); i++ {
		product := productsReturn[i]

		if priceLessThan == nil || product.Price < *priceLessThan {
			filteredProducts = append(filteredProducts, *product)
		}

		if len(filteredProducts) == 5 {
			break
		}
	}
	var nextCursor = ""
	if len(filteredProducts) == 5 && len(productsReturn) > startIndex+len(filteredProducts) {
		nextCursor = strconv.Itoa(startIndex + len(filteredProducts))
	}

	return filteredProducts, nextCursor, nil
}
