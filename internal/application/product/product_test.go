package product_application_test

import (
	product_application "mytheresa/internal/application/product"
	"mytheresa/internal/models"
	"mytheresa/internal/storage/mock"
	"testing"
)

func TestGetProducts_ExpectedProducts(t *testing.T) {
	var mockProductStorage mock.MockProduct
	var mockDiscountStorage mock.MockDiscount

	err := mockProductStorage.LoadData()
	if err != nil {
		t.Fatalf("Failed to load mock product data: %v", err)
	}

	err = mockDiscountStorage.LoadData()
	if err != nil {
		t.Fatalf("Failed to load mock discount data: %v", err)
	}

	dProduct := &product_application.ProductApplication{
		StgProduct:  &mockProductStorage,
		StgDiscount: &mockDiscountStorage,
	}

	category := ""
	cursor := ""

	productsPrice, nextCursor, err := dProduct.GetProducts(category, nil, cursor)
	if err != nil {
		t.Fatalf("GetProducts() returned an error: %v", err)
	}

	expectedResults := []models.ProductPrice{
		{
			SKU:      "000002",
			Name:     "Product C",
			Category: "boots",
			Price: models.Price{
				Original: 50000,
				Final:    35000,
				Discount: ptr("30%"),
				Currency: "EUR",
			},
		},
		{
			SKU:      "000005",
			Name:     "Product E",
			Category: "sandals",
			Price: models.Price{
				Original: 62000,
				Final:    62000,
				Discount: nil,
				Currency: "EUR",
			},
		},
		{
			SKU:      "000016",
			Name:     "Product P",
			Category: "boots",
			Price: models.Price{
				Original: 64000,
				Final:    44800,
				Discount: ptr("30%"),
				Currency: "EUR",
			},
		},
		{
			SKU:      "000011",
			Name:     "Product K",
			Category: "boots",
			Price: models.Price{
				Original: 65000,
				Final:    45500,
				Discount: ptr("30%"),
				Currency: "EUR",
			},
		},
		{
			SKU:      "000003",
			Name:     "Product B",
			Category: "sandals",
			Price: models.Price{
				Original: 66000,
				Final:    56100,
				Discount: ptr("15%"),
				Currency: "EUR",
			},
		},
	}

	if len(productsPrice) != len(expectedResults) {
		t.Fatalf("Expected %d products in the result, got %d", len(expectedResults), len(productsPrice))
	}

	for i, productPrice := range productsPrice {
		expected := expectedResults[i]

		if productPrice.SKU != expected.SKU {
			t.Errorf("At index %d, expected SKU %s, got %s", i, expected.SKU, productPrice.SKU)
		}

		if productPrice.Name != expected.Name {
			t.Errorf("At index %d, expected Name %s, got %s", i, expected.Name, productPrice.Name)
		}

		if productPrice.Category != expected.Category {
			t.Errorf("At index %d, expected Category %s, got %s", i, expected.Category, productPrice.Category)
		}

		if productPrice.Price.Original != expected.Price.Original {
			t.Errorf("At index %d, expected Original Price %d, got %d", i, expected.Price.Original, productPrice.Price.Original)
		}

		if productPrice.Price.Final != expected.Price.Final {
			t.Errorf("At index %d, expected Final Price %d, got %d", i, expected.Price.Final, productPrice.Price.Final)
		}

		if !safeCompareDiscount(productPrice.Price.Discount, expected.Price.Discount) {
			t.Errorf("At index %d, expected Discount %v, got %v", i, expected.Price.Discount, productPrice.Price.Discount)
		}

		if productPrice.Price.Currency != expected.Price.Currency {
			t.Errorf("At index %d, expected Currency %s, got %s", i, expected.Price.Currency, productPrice.Price.Currency)
		}
	}

	if nextCursor != "5" {
		t.Errorf("Expected nextCursor to be empty, got %s", nextCursor)
	}
}

func ptr(s string) *string {
	return &s
}

func safeCompareDiscount(actual, expected *string) bool {
	if actual == nil && expected == nil {
		return true
	}
	if actual == nil || expected == nil {
		return false
	}
	return *actual == *expected
}
