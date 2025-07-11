package main

import "time"

// Domain models
type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Order struct {
	ID         string      `json:"id"`
	CustomerID string      `json:"customer_id"`
	Items      []OrderItem `json:"items"`
	Status     string      `json:"status"` // "pending", "confirmed", "cancelled"
	Total      float64     `json:"total"`
	CreatedAt  time.Time   `json:"created_at"`
}

// Request/Response DTOs
type CreateOrderRequest struct {
	CustomerID string                   `json:"customer_id"`
	Items      []CreateOrderRequestItem `json:"items"`
}

type CreateOrderRequestItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
