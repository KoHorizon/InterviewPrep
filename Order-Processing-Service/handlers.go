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

// POST /orders - Create a new order
func (h *Handlers) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate customer ID
	if req.CustomerID == "" {
		h.respondWithError(w, "customer_id is required", http.StatusBadRequest)
		return
	}

	// Validate items exist
	if len(req.Items) == 0 {
		h.respondWithError(w, "at least one item is required", http.StatusBadRequest)
		return
	}

	// Build order items with current prices
	orderItems := make([]OrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		// Validate product ID
		if item.ProductID == "" {
			h.respondWithError(w, "product_id is required for all items", http.StatusBadRequest)
			return
		}

		// Validate quantity
		if item.Quantity <= 0 {
			h.respondWithError(w, "quantity must be greater than 0", http.StatusBadRequest)
			return
		}

		// Get product to check existence and get current price
		product, exists := h.store.GetProduct(item.ProductID)
		if !exists {
			h.respondWithError(w, "product not found: "+item.ProductID, http.StatusBadRequest)
			return
		}

		orderItems = append(orderItems, OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price, // Use current price
		})
	}

	// Create order (this will validate stock and deduct inventory)
	order, err := h.store.CreateOrder(req.CustomerID, orderItems)
	if err != nil {
		h.respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return success response
	h.respondWithJSON(w, order, http.StatusCreated)
}

// GET /orders/{id} - Get order by ID
func (h *Handlers) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	order, exists := h.store.GetOrder(orderID)
	if !exists {
		h.respondWithError(w, "order not found", http.StatusNotFound)
		return
	}

	h.respondWithJSON(w, order, http.StatusOK)
}

// PUT /orders/{id}/cancel - Cancel an order
func (h *Handlers) CancelOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	err := h.store.CancelOrder(orderID)
	if err != nil {
		// Determine appropriate status code
		status := http.StatusBadRequest
		if err.Error() == "order not found: "+orderID {
			status = http.StatusNotFound
		}
		h.respondWithError(w, err.Error(), status)
		return
	}

	// Get updated order to return
	order, _ := h.store.GetOrder(orderID)
	h.respondWithJSON(w, order, http.StatusOK)
}

// Helper methods
func (h *Handlers) respondWithJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handlers) respondWithError(w http.ResponseWriter, message string, status int) {
	h.respondWithJSON(w, ErrorResponse{Error: message}, status)
}
