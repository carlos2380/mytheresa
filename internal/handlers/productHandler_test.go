package handlers_test

import (
	"encoding/json"
	product_application "mytheresa/internal/application/product"
	"mytheresa/internal/handlers"
	"mytheresa/internal/storage/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProductsHandler(t *testing.T) {

	mockProductStorage := &mock.MockProduct{}
	mockDiscountStorage := &mock.MockDiscount{}

	err := mockProductStorage.LoadData()
	if err != nil {
		t.Fatalf("Failed to load mock product data: %v", err)
	}

	err = mockDiscountStorage.LoadData()
	if err != nil {
		t.Fatalf("Failed to load mock discount data: %v", err)
	}

	productApp := &product_application.ProductApplication{
		StgProduct:  mockProductStorage,
		StgDiscount: mockDiscountStorage,
	}

	handler := handlers.ProductHandler{
		AppProduct: productApp,
	}

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "Valid request with priceLessThan filter",
			method:         http.MethodGet,
			url:            "/products?priceLessThan=64000",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"products": []map[string]interface{}{
					{
						"sku":      "000002",
						"name":     "Product C",
						"category": "boots",
						"price": map[string]interface{}{
							"original":            50000,
							"value":               35000,
							"discount_percentage": "30%",
							"currency":            "EUR",
						},
					},
					{
						"sku":      "000005",
						"name":     "Product E",
						"category": "sandals",
						"price": map[string]interface{}{
							"original":            62000,
							"value":               62000,
							"discount_percentage": nil,
							"currency":            "EUR",
						},
					},
				},
				"nextCursor": "",
			},
		},
		{
			name:           "Invalid price filter",
			method:         http.MethodGet,
			url:            "/products?priceLessThan=abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":  400,
				"error": "Invalid priceLessThan parameter: must be an integer",
			},
		},
		{
			name:           "Method not allowed (POST)",
			method:         http.MethodPost,
			url:            "/products",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.url, nil)
			rec := httptest.NewRecorder()

			handler.GetProducts(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, res.StatusCode)
			}

			if tc.expectedBody != nil {
				var actualBody map[string]interface{}
				err := json.NewDecoder(res.Body).Decode(&actualBody)
				if err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				if !deepEqualJSON(tc.expectedBody, actualBody) {
					t.Errorf("Unexpected response body.\nGot: %+v\nWant: %+v", actualBody, tc.expectedBody)
				}
			} else {
				body := rec.Body.String()
				if body != "" {
					t.Errorf("Expected empty body, got: %s", body)
				}
			}
		})
	}
}

func deepEqualJSON(expected, actual interface{}) bool {
	expectedJSON, _ := json.Marshal(expected)
	actualJSON, _ := json.Marshal(actual)
	return string(expectedJSON) == string(actualJSON)
}
