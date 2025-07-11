package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize store with sample data
	store := NewStore()
	initializeProducts(store)

	// Initialize handlers
	handlers := NewHandlers(store)

	// Setup routes
	router := mux.NewRouter()
	router.HandleFunc("/orders", handlers.CreateOrder).Methods("POST")
	router.HandleFunc("/orders/{id}", handlers.GetOrder).Methods("GET")
	router.HandleFunc("/orders/{id}/cancel", handlers.CancelOrder).Methods("PUT")

	// Start server
	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func initializeProducts(store *Store) {
	products := []Product{
		{ID: "prod1", Name: "Laptop", Price: 999.99, Stock: 5},
		{ID: "prod2", Name: "Mouse", Price: 29.99, Stock: 10},
		{ID: "prod3", Name: "Keyboard", Price: 79.99, Stock: 3},
	}

	for _, product := range products {
		store.AddProduct(product)
	}
}
