package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Handlers struct {
	store *Store
}

func NewHandlers(store *Store) *Handlers {
	return &Handlers{store: store}
}

func (h *Handlers) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.CustomerID == "" {
		h.sendError(w, "customer_id is required", http.StatusBadRequest)
		return
	}
	if len(req.Items) == 0 {
		h.sendError(w, "items are required", http.StatusBadRequest)
		return
	}

	// Convert request items to order items with current prices
	var orderItems []OrderItem
	for _, item := range req.Items {
		if item.Quantity <= 0 {
			h.sendError(w, "quantity must be positive", http.StatusBadRequest)
			return
		}

		product, exists := h.store.GetProduct(item.ProductID)
		if !exists {
			h.sendError(w, "product not found: "+item.ProductID, http.StatusBadRequest)
			return
		}

		orderItems = append(orderItems, OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
	}

	// Create order
	order, err := h.store.CreateOrder(req.CustomerID, orderItems)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *Handlers) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	order, exists := h.store.GetOrder(orderID)
	if !exists {
		h.sendError(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *Handlers) CancelOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	err := h.store.CancelOrder(orderID)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "order "+orderID+" not found" {
			status = http.StatusNotFound
		}
		h.sendError(w, err.Error(), status)
		return
	}

	// Return updated order
	order, _ := h.store.GetOrder(orderID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (h *Handlers) sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
