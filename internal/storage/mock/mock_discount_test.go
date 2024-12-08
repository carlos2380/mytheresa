package mock_test

import (
	"fmt"
	"mytheresa/internal/storage"
	"mytheresa/internal/storage/mock"
	"testing"
)

func NewMockDiscountWithData() *mock.MockDiscount {
	mockDiscount := &mock.MockDiscount{}
	err := mockDiscount.LoadData()
	if err != nil {
		panic(fmt.Sprintf("failed to load mock data: %v", err))
	}
	return mockDiscount
}

func TestLoadData(t *testing.T) {

	var mockDiscount storage.DiscountStoreage = NewMockDiscountWithData()

	category := "boots"
	percent, _ := mockDiscount.GetDiscount(category)
	if percent != 30 {
		t.Errorf("Expected discount percentage for category %s to be 30, got %d", category, percent)
	}

	nonExistentCategory := "nonExistentCategory"
	percent, _ = mockDiscount.GetDiscount(nonExistentCategory)
	if percent != 0 {
		t.Errorf("Expected 0 for non-existent category, got %d", percent)
	}
}
