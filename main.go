package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/seyyidibrahimgulec/product-listing/app/handlers"
)

func main() {
	r := mux.NewRouter()
	handlers.RegisterProductRoutes(r)
	http.Handle("/", r)
	log.Println("Listening and serving on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
