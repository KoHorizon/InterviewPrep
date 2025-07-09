// store.go
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
		products:     make(map[string]Product),
		orders:       make(map[string]Order),
		orderCounter: 0,
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

	// Validate inventory and calculate total
	var total float64
	for _, item := range items {
		product, exists := s.products[item.ProductID]
		if !exists {
			return Order{}, fmt.Errorf("product %s not found", item.ProductID)
		}
		if product.Stock < item.Quantity {
			return Order{}, fmt.Errorf("insufficient stock for product %s", item.ProductID)
		}
		total += item.Price * float64(item.Quantity)
	}

	// Generate order ID
	s.orderCounter++
	orderID := fmt.Sprintf("order_%d", s.orderCounter)

	// Deduct inventory
	for _, item := range items {
		product := s.products[item.ProductID]
		product.Stock -= item.Quantity
		s.products[item.ProductID] = product
	}

	// Create order
	order := Order{
		ID:         orderID,
		CustomerID: customerID,
		Items:      items,
		Status:     "pending",
		Total:      total,
		CreatedAt:  time.Now(),
	}

	s.orders[orderID] = order
	return order, nil
}

func (s *Store) CancelOrder(orderID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[orderID]
	if !exists {
		return fmt.Errorf("order %s not found", orderID)
	}

	if order.Status != "pending" {
		return fmt.Errorf("order %s cannot be cancelled, status: %s", orderID, order.Status)
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
