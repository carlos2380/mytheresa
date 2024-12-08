package mock

import (
	"encoding/json"
	"log"
	"mytheresa/internal/models"
)

type MockDiscount struct {
	Discounts map[string]models.DiscountCategory // Cambiado a un mapa
}

func (mDiscount *MockDiscount) LoadData() error {
	if mDiscount.Discounts == nil {
		mDiscount.Discounts = make(map[string]models.DiscountCategory)
	}

	var discountList []models.DiscountCategory
	err := json.Unmarshal([]byte(rawDiscountsJSON), &discountList)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		return err
	}

	for _, discount := range discountList {
		mDiscount.Discounts[discount.Category] = discount
	}

	return nil
}

func (mDiscount *MockDiscount) GetDiscount(category string) (int, error) {
	discount, exists := mDiscount.Discounts[category]
	if !exists {
		return 0, nil
	}
	return discount.Percent, nil
}
