package mock_test

import (
	"fmt"
	"mytheresa/internal/models"
	"mytheresa/internal/storage"
	"mytheresa/internal/storage/mock"
	"testing"
)

func NewMockProductWithData() *mock.MockProduct {
	mockProduct := &mock.MockProduct{}
	err := mockProduct.LoadData()
	if err != nil {
		panic(fmt.Sprintf("failed to load mock data: %v", err))
	}
	return mockProduct
}

func TestGetProduct(t *testing.T) {
	var mockProduct storage.ProductStoreage = NewMockProductWithData()

	tests := []struct {
		name               string
		category           string
		priceLessThan      *int
		cursor             string
		expectedCount      int
		expectedNextCursor string
		expectedError      bool
		expectedSKUs       []string
	}{
		{
			name:               "Fetch first 5 products from 'boots' category",
			category:           "boots",
			priceLessThan:      nil,
			cursor:             "",
			expectedCount:      5,
			expectedNextCursor: "5",
			expectedError:      false,
			expectedSKUs:       []string{"000002", "000016", "000011", "000019", "000009"},
		},
		{
			name:               "Continue fetching products with cursor",
			category:           "boots",
			priceLessThan:      nil,
			cursor:             "5",
			expectedCount:      3,
			expectedNextCursor: "",
			expectedError:      false,
			expectedSKUs:       []string{"000006", "000013", "000001"},
		},
		{
			name:               "Fetch products with price less than 70,000 in 'boots' category",
			category:           "boots",
			priceLessThan:      ptrInt(70000),
			cursor:             "",
			expectedCount:      3,
			expectedNextCursor: "",
			expectedError:      false,
			expectedSKUs:       []string{"000002", "000016", "000011"},
		},
		{
			name:               "Fetch products from a non-existent category",
			category:           "non-existing-category",
			priceLessThan:      nil,
			cursor:             "",
			expectedCount:      0,
			expectedNextCursor: "",
			expectedError:      false,
			expectedSKUs:       []string{},
		},
		{
			name:               "Invalid cursor",
			category:           "boots",
			priceLessThan:      nil,
			cursor:             "invalid",
			expectedCount:      0,
			expectedNextCursor: "",
			expectedError:      true,
			expectedSKUs:       nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			products, nextCursor, err := mockProduct.GetProducts(test.category, test.priceLessThan, test.cursor)

			if test.expectedError && err == nil {
				t.Errorf("Expected an error but got none")
			} else if !test.expectedError && err != nil {
				t.Errorf("Did not expect an error but got: %v", err)
			}

			if len(products) != test.expectedCount {
				t.Errorf("Expected %d products but got %d", test.expectedCount, len(products))
			}

			if nextCursor != test.expectedNextCursor {
				t.Errorf("Expected next cursor to be '%s', but got '%s'", test.expectedNextCursor, nextCursor)
			}

			if test.expectedSKUs != nil {
				actualSKUs := extractSKUs(products)
				if !equalSlices(actualSKUs, test.expectedSKUs) {
					t.Errorf("Expected SKUs %v but got %v", test.expectedSKUs, actualSKUs)
				}
			}
		})
	}
}

func ptrInt(i int) *int {
	return &i
}

func extractSKUs(products []models.Product) []string {
	skus := make([]string, len(products))
	for i, product := range products {
		skus[i] = product.SKU
	}
	return skus
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
