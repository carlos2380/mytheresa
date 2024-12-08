package server

import (
	"mytheresa/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(pHandler *handlers.ProductHandler) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/api/products", pHandler.GetProducts).Methods(http.MethodGet, http.MethodOptions)
	return r
}
