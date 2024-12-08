package handlers

import (
	"encoding/json"
	"log"
	"mytheresa/internal/errors"
	"net/http"
	"strconv"
)

func (pHandler ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodGet:
		category := r.URL.Query().Get("category")
		priceLessThan := r.URL.Query().Get("priceLessThan")
		cursor := r.URL.Query().Get("cursor")
		var priceLessThanInt *int
		if priceLessThan != "" {
			priceLess, err := strconv.Atoi(priceLessThan)
			if err != nil {
				log.Println(errors.Wrap(err, *errors.ErrPriceLessConvert))
				errors.ErrPriceLessConvert.Respond(w)
				return
			}
			priceLessThanInt = &priceLess
		}

		productPrice, nextCursor, err := pHandler.AppProduct.GetProducts(category, priceLessThanInt, cursor)
		if err != nil {
			log.Printf("Error encoding JSON response: %v", err)
			errors.ErrInternalServerError.Respond(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"products":   productPrice,
			"nextCursor": nextCursor,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding JSON response: %v", err)
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
