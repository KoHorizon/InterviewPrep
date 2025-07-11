package main

import (
	"fmt"
	"sync"
	"time"
)

type Store struct {
	mu           sync.RWMutex
	products     map[string]Product
	orders       map[string]Order
	orderCounter int
}

func NewStore() *Store {
	return &Store{
		products: make(map[string]Product),
		orders:   make(map[string]Order),
	}
}

func (s *Store) AddProduct(product Product) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.products[product.ID] = product
}

func (s *Store) GetProduct(id string) (Product, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	product, exists := s.products[id]
	return product, exists
}

func (s *Store) GetOrder(id string) (Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, exists := s.orders[id]
	return order, exists
}

func (s *Store) CreateOrder(customerID string, items []OrderItem) (Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// First, validate all items have sufficient stock
	for _, item := range items {
		product, exists := s.products[item.ProductID]
		if !exists {
			return Order{}, fmt.Errorf("product not found: %s", item.ProductID)
		}
		if product.Stock < item.Quantity {
			return Order{}, fmt.Errorf("insufficient stock for product %s: available %d, requested %d",
				item.ProductID, product.Stock, item.Quantity)
		}
	}

	// Calculate total
	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}

	// Generate order ID
	s.orderCounter++
	orderID := fmt.Sprintf("order_%d", s.orderCounter)

	// Create order
	order := Order{
		ID:         orderID,
		CustomerID: customerID,
		Items:      items,
		Status:     "pending",
		Total:      total,
		CreatedAt:  time.Now(),
	}

	// Deduct inventory only after order is created
	for _, item := range items {
		product := s.products[item.ProductID]
		product.Stock -= item.Quantity
		s.products[item.ProductID] = product
	}

	// Save order
	s.orders[orderID] = order

	return order, nil
}

func (s *Store) CancelOrder(orderID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if order exists
	order, exists := s.orders[orderID]
	if !exists {
		return fmt.Errorf("order not found: %s", orderID)
	}

	// Check if order can be cancelled
	if order.Status != "pending" {
		return fmt.Errorf("cannot cancel order with status: %s", order.Status)
	}

	// Restore inventory
	for _, item := range order.Items {
		product := s.products[item.ProductID]
		product.Stock += item.Quantity
		s.products[item.ProductID] = product
	}

	// Update order status
	order.Status = "cancelled"
	s.orders[orderID] = order

	return nil
}
