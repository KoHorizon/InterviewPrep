#!/bin/bash

# Core features test script

echo "=== Testing Core Features Only ==="
echo ""

# 1. Create a valid order
echo "1. Creating a valid order:"
curl -s -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "cust123",
    "items": [
      {"product_id": "prod1", "quantity": 1},
      {"product_id": "prod2", "quantity": 2}
    ]
  }' | jq .
echo ""

# 2. Get the created order
echo "2. Getting order details:"
curl -s http://localhost:8080/orders/order_1 | jq .
echo ""

# 3. Try to create order with insufficient stock
echo "3. Testing insufficient stock (should fail):"
curl -s -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "cust456",
    "items": [
      {"product_id": "prod3", "quantity": 5}
    ]
  }' | jq .
echo ""

# 4. Cancel the order
echo "4. Cancelling order:"
curl -s -X PUT http://localhost:8080/orders/order_1/cancel | jq .
echo ""

# 5. Try to cancel already cancelled order
echo "5. Trying to cancel already cancelled order (should fail):"
curl -s -X PUT http://localhost:8080/orders/order_1/cancel | jq .
echo ""

# 6. Get non-existent order
echo "6. Getting non-existent order (should return 404):"
curl -s http://localhost:8080/orders/order_999 | jq .
echo ""

# 7. Test validation - no customer_id
echo "7. Testing validation - missing customer_id:"
curl -s -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "items": [{"product_id": "prod1", "quantity": 1}]
  }' | jq .
echo ""

# 8. Test validation - empty items
echo "8. Testing validation - empty items:"
curl -s -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "cust123",
    "items": []
  }' | jq .
echo ""

# 9. Test validation - invalid quantity
echo "9. Testing validation - invalid quantity:"
curl -s -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "cust123",
    "items": [{"product_id": "prod1", "quantity": 0}]
  }' | jq .
echo ""

# 10. Test validation - non-existent product
echo "10. Testing validation - non-existent product:"
curl -s -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "cust123",
    "items": [{"product_id": "prod999", "quantity": 1}]
  }' | jq .

echo ""
echo "=== Core Features Test Complete ==="
