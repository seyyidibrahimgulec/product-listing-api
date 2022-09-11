package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/seyyidibrahimgulec/product-listing/app/models"
)

func logError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func getProductRespondError(w http.ResponseWriter, r *http.Request) *models.Product {
	vars := mux.Vars(r)
	productId := vars["productId"]
	id, err := strconv.ParseInt(productId, 0, 0)
	if err != nil {
		logError(err)
		createJSONResponse(w, http.StatusBadRequest, []byte{})
		return nil
	}
	product, err := models.GetProductById(id)
	if err != nil {
		logError(err)
		createJSONResponse(w, http.StatusNotFound, []byte{})
		w.WriteHeader(http.StatusNotFound)
		return nil
	}
	return product
}

func createJSONResponse(w http.ResponseWriter, status_code int, response []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)
	w.Write(response)
}

func listProducts(w http.ResponseWriter, r *http.Request) {
	products := models.GetAllProducts()
	res, err := json.Marshal(products)
	logError(err)
	createJSONResponse(w, http.StatusOK, res)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	product := models.Product{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&product); err != nil {
		createJSONResponse(w, http.StatusBadRequest, []byte(err.Error()))
		return
	}
	defer r.Body.Close()

	p, err := product.CreateProduct()

	if err != nil {
		createJSONResponse(w, http.StatusInternalServerError, []byte{})
		log.Println([]byte(err.Error()))
	} else {
		res, err := json.Marshal(p)
		logError(err)
		createJSONResponse(w, http.StatusCreated, res)
	}
}

func retrieveProduct(w http.ResponseWriter, r *http.Request) {
	product := getProductRespondError(w, r)
	if product == nil {
		return
	}
	res, err := json.Marshal(product)
	if err != nil {
		logError(err)
		createJSONResponse(w, http.StatusInternalServerError, []byte{})
		return
	}
	createJSONResponse(w, http.StatusOK, res)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	product := getProductRespondError(w, r)
	if product == nil {
		return
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&product); err != nil {
		createJSONResponse(w, http.StatusBadRequest, []byte(err.Error()))
		return
	}
	defer r.Body.Close()

	product, err := product.UpdateProduct()
	logError(err)
	res, err := json.Marshal(product)
	logError(err)
	createJSONResponse(w, http.StatusOK, res)
}

var RegisterProductRoutes = func(router *mux.Router) {
	router.HandleFunc("/product/", listProducts).Methods("GET")
	router.HandleFunc("/product/", createProduct).Methods("POST")
	router.HandleFunc("/product/{productId}/", retrieveProduct).Methods("GET")
	router.HandleFunc("/product/{productId}/", updateProduct).Methods("PUT")
}
