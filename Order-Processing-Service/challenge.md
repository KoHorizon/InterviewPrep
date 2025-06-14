# Go Backend Challenge: Order Processing Service

## Challenge Overview
**Duration**: 60 minutes
**Role**: Mid-level Go Backend Developer
**Focus**: Practical backend skills, REST API development, concurrency handling

---

## Problem Statement

You're building a microservice for an e-commerce platform that processes customer orders. The service needs to handle order creation, inventory validation, and basic order management.

## Data Models

```go
type Product struct {
    ID    string
    Name  string
    Price float64
    Stock int
}

type OrderItem struct {
    ProductID string
    Quantity  int
    Price     float64
}

type Order struct {
    ID          string
    CustomerID  string
    Items       []OrderItem
    Status      string // "pending", "confirmed", "cancelled"
    Total       float64
    CreatedAt   time.Time
}
```

## Requirements

### Core Features (45 minutes)

#### 1. POST /orders - Create a new order
- Validate inventory availability for all items
- Calculate total price from current product prices
- Return 400 Bad Request if insufficient stock for any item
- Return 201 Created with order details on success
- Deduct inventory when order is created

#### 2. GET /orders/{id} - Retrieve order by ID
- Return 404 Not Found if order doesn't exist
- Return 200 OK with complete order details

#### 3. PUT /orders/{id}/cancel - Cancel an order
- Only allow cancellation if current status is "pending"
- Restore inventory quantities when order is cancelled
- Return 400 Bad Request if order cannot be cancelled
- Return 200 OK when successfully cancelled

### Technical Requirements

- **Storage**: Use in-memory storage (maps) - no database required
- **HTTP Handling**: Proper HTTP status codes and error responses
- **Validation**: Basic input validation for all endpoints
- **Concurrency**: Thread-safe operations (handle concurrent requests)
- **Code Quality**: Clean, readable code with proper structure
- **Error Handling**: Comprehensive error handling with meaningful messages

### Bonus Features (if time permits)

- **GET /orders** - List all orders with optional status filter (`?status=pending`)
- **Logging**: Basic request/response logging
- **Middleware**: Simple middleware (e.g., request logging, CORS)
- **Health Check**: GET /health endpoint

## Sample Test Data

Pre-populate your service with these products:

```go
products := map[string]Product{
    "prod1": {ID: "prod1", Name: "Laptop", Price: 999.99, Stock: 5},
    "prod2": {ID: "prod2", Name: "Mouse", Price: 29.99, Stock: 10},
    "prod3": {ID: "prod3", Name: "Keyboard", Price: 79.99, Stock: 3},
}
```

## Sample API Requests

### Create Order
```bash
POST /orders
Content-Type: application/json

{
    "customer_id": "cust123",
    "items": [
        {
            "product_id": "prod1",
            "quantity": 1
        },
        {
            "product_id": "prod2",
            "quantity": 2
        }
    ]
}
```

### Get Order
```bash
GET /orders/order123
```

### Cancel Order
```bash
PUT /orders/order123/cancel
```

## Expected Deliverables

1. **Working HTTP Server**: Complete Go application with all required endpoints
2. **Proper Error Handling**: Appropriate HTTP status codes and error messages
3. **Thread Safety**: Safe concurrent access to shared data
4. **Clean Architecture**: Well-structured code with separated concerns
5. **Input Validation**: Proper validation of request data
6. **Documentation**: Brief README or comments explaining your approach

## Evaluation Criteria

### Technical Skills (70%)
- **Code Structure**: Clean separation of handlers, business logic, and data access
- **Concurrency Safety**: Proper use of mutexes, channels, or other synchronization
- **Error Handling**: Comprehensive error handling with appropriate HTTP responses
- **HTTP Best Practices**: Correct use of HTTP methods, status codes, and JSON
- **Input Validation**: Thorough validation of request bodies and parameters
- **Go Idioms**: Idiomatic Go code following community standards

### Problem-Solving (20%)
- Approach to breaking down the problem
- Handling of edge cases and error scenarios
- Questions asked for clarification
- Solutions to concurrency challenges

### Communication (10%)
- Clear explanation of thought process
- Relevant questions during development
- Code readability and appropriate comments

## Follow-up Discussion Topics

### During Implementation
- "How would you handle two customers ordering the last item simultaneously?"
- "What should happen if someone tries to cancel a confirmed order?"
- "How would you structure this to add payment processing later?"

### After Implementation
- "How would you adapt this to work with a PostgreSQL database?"
- "What would you add to make this production-ready?"
- "How would you handle high traffic scenarios?"
- "What testing strategy would you implement?"
- "How would you deploy and monitor this service in production?"
- "What metrics would you track for this service?"

## Success Indicators

### Green Flags ✅
- Clean, idiomatic Go code
- Proper concurrency handling
- Thoughtful API design
- Comprehensive error handling
- Good separation of concerns
- Questions about edge cases and production considerations

### Red Flags ❌
- Race conditions or unsafe concurrent access
- Poor or missing error handling
- Overly complex solutions
- No input validation
- Hardcoded configuration values
- Ignoring HTTP best practices

## Time Management Tips

- **0-15 minutes**: Set up basic server structure and routing
- **15-30 minutes**: Implement core data structures and business logic
- **30-45 minutes**: Complete all three main endpoints
- **45-60 minutes**: Add error handling, validation, and polish

## Getting Started

1. Initialize a new Go module
2. Set up basic HTTP server with routing
3. Define your data structures
4. Implement in-memory storage with proper synchronization
5. Build endpoints incrementally
6. Test with sample requests
7. Add error handling and validation

Good luck! Focus on building a working solution first, then improve code quality and add features.
