package main

import (
	"log"
	"net/http"
	"projectGolang/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Product routes
	r.HandleFunc("/products", handlers.GetProducts).Methods("GET")
	r.HandleFunc("/products/{id}", handlers.GetProduct).Methods("GET")
	r.HandleFunc("/products", handlers.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", handlers.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", handlers.DeleteProduct).Methods("DELETE")

	// Category routes
	r.HandleFunc("/categories", handlers.GetCategories).Methods("GET")
	r.HandleFunc("/categories", handlers.CreateCategory).Methods("POST")

	// User routes
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")

	log.Println("Server running on port 8000")
	http.ListenAndServe(":8000", r)
}
