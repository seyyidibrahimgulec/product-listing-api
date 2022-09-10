package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/seyyidibrahimgulec/product-listing/app/models"
)

func listProducts(w http.ResponseWriter, r *http.Request) {
	products := models.GetAllProducts()
	res, _ := json.Marshal(products)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	product := models.Product{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	p, err := product.CreateProduct()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		res, _ := json.Marshal(p)
		w.WriteHeader(http.StatusCreated)
		w.Write(res)
	}
}

var RegisterProductRoutes = func(router *mux.Router) {
	router.HandleFunc("/product/", listProducts).Methods("GET")
	router.HandleFunc("/product/", createProduct).Methods("POST")
}
