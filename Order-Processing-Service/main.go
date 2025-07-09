// main.go
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

	// Initialize handlers with store dependency
	h := NewHandlers(store)

	// Setup routes
	r := mux.NewRouter()
	r.HandleFunc("/orders", h.CreateOrder).Methods("POST")
	r.HandleFunc("/orders/{id}", h.GetOrder).Methods("GET")
	r.HandleFunc("/orders/{id}/cancel", h.CancelOrder).Methods("PUT")
	r.HandleFunc("/health", h.Health).Methods("GET")

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func initializeProducts(store *Store) {
	products := map[string]Product{
		"prod1": {ID: "prod1", Name: "Laptop", Price: 999.99, Stock: 5},
		"prod2": {ID: "prod2", Name: "Mouse", Price: 29.99, Stock: 10},
		"prod3": {ID: "prod3", Name: "Keyboard", Price: 79.99, Stock: 3},
	}

	for _, product := range products {
		store.AddProduct(product)
	}
}
