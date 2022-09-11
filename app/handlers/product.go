package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	w.Header().Set("Content-Type", "application/json")
	if err := decoder.Decode(&product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	p, err := product.CreateProduct()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println([]byte(err.Error()))
	} else {
		res, _ := json.Marshal(p)
		w.WriteHeader(http.StatusCreated)
		w.Write(res)
	}
}

func retrieveProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	productId := vars["productId"]
	id, err := strconv.ParseInt(productId, 0, 0)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	productDetails, err := models.GetProductById(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	res, err := json.Marshal(productDetails)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["productId"]
	id, err := strconv.ParseInt(productId, 0, 0)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product, err := models.GetProductById(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	decoder := json.NewDecoder(r.Body)
	w.Header().Set("Content-Type", "application/json")
	if err := decoder.Decode(&product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	product, err = product.UpdateProduct()
	if err != nil {
		log.Println(err)
	}
	res, _ := json.Marshal(product)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

var RegisterProductRoutes = func(router *mux.Router) {
	router.HandleFunc("/product/", listProducts).Methods("GET")
	router.HandleFunc("/product/", createProduct).Methods("POST")
	router.HandleFunc("/product/{productId}/", retrieveProduct).Methods("GET")
	router.HandleFunc("/product/{productId}/", updateProduct).Methods("PUT")
}
