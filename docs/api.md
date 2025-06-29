# Transaction Management API Documentation

## Overview
RESTful API for transaction management with dashboard analytics built using Go, Gin, GORM, and MySQL.

## Base URL
```
http://localhost:8080/api
```

## Authentication
This API does not require authentication in the current implementation.

## Response Format
All responses follow this standard format:

```json
{
  "success": boolean,
  "data": object|array|null,
  "message": string,
  "error": string
}
```

## Endpoints

### 1. Create Transaction
**POST** `/transactions`

Creates a new transaction with pending status.

**Request Body:**
```json
{
  "user_id": 1,
  "amount": 100.50
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "amount": 100.50,
    "status": "pending",
    "created_at": "2025-06-28T10:00:00Z",
    "updated_at": "2025-06-28T10:00:00Z"
  },
  "message": "Transaction created successfully"
}
```

### 2. Get All Transactions
**GET** `/transactions`

Retrieves all transactions with optional filtering and pagination.

**Query Parameters:**
- `user_id` (integer, optional): Filter by user ID
- `status` (string, optional): Filter by status (pending, success, failed)
- `limit` (integer, optional): Number of records to return (default: 20, max: 100)
- `offset` (integer, optional): Number of records to skip (default: 0)

**Example:**
```
GET /transactions?user_id=1&status=pending&limit=10&offset=0
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "amount": 100.50,
      "status": "pending",
      "created_at": "2025-06-28T10:00:00Z",
      "updated_at": "2025-06-28T10:00:00Z"
    }
  ],
  "message": "Transactions retrieved successfully"
}
```

### 3. Get Transaction by ID
**GET** `/transactions/{id}`

Retrieves a specific transaction by its ID.

**Path Parameters:**
- `id` (integer): Transaction ID

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "amount": 100.50,
    "status": "pending",
    "created_at": "2025-06-28T10:00:00Z",
    "updated_at": "2025-06-28T10:00:00Z"
  },
  "message": "Transaction retrieved successfully"
}
```

### 4. Update Transaction Status
**PUT** `/transactions/{id}`

Updates the status of a specific transaction.

**Path Parameters:**
- `id` (integer): Transaction ID

**Request Body:**
```json
{
  "status": "success"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": null,
  "message": "Transaction updated successfully"
}
```

### 5. Delete Transaction
**DELETE** `/transactions/{id}`

Deletes a specific transaction.

**Path Parameters:**
- `id` (integer): Transaction ID

**Response (200 OK):**
```json
{
  "success": true,
  "data": null,
  "message": "Transaction deleted successfully"
}
```

### 6. Dashboard Summary
**GET** `/dashboard/summary`

Retrieves dashboard summary with analytics data.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "today_successful_transactions": 5,
    "today_successful_amount": 1250.75,
    "average_transaction_per_user": 3.2,
    "latest_transactions": [
      {
        "id": 10,
        "user_id": 2,
        "amount": 250.00,
        "status": "success",
        "created_at": "2025-06-28T14:30:00Z",
        "updated_at": "2025-06-28T14:35:00Z"
      }
    ],
    "status_counts": {
      "success": 15,
      "pending": 8,
      "failed": 2
    }
  },
  "message": "Dashboard summary retrieved successfully"
}
```

## Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "error": "Invalid request body"
}
```

### 404 Not Found
```json
{
  "success": false,
  "error": "Transaction not found"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "error": "Internal server error"
}
```

## Status Codes
- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Validation Rules

### Create Transaction
- `user_id`: Required, must be positive integer
- `amount`: Required, must be positive number (minimum 0.01)

### Update Transaction
- `status`: Required, must be one of: "pending", "success", "failed"

## Health Check
**GET** `/health`

Returns server health status.

**Response (200 OK):**
```json
{
  "status": "OK"
}
```

## Testing with cURL

### Create a transaction:
```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "amount": 100.50}'
```

### Get all transactions:
```bash
curl http://localhost:8080/api/transactions
```

### Get transaction by ID:
```bash
curl http://localhost:8080/api/transactions/1
```

### Update transaction status:
```bash
curl -X PUT http://localhost:8080/api/transactions/1 \
  -H "Content-Type: application/json" \
  -d '{"status": "success"}'
```

### Delete transaction:
```bash
curl -X DELETE http://localhost:8080/api/transactions/1
```

### Get dashboard summary:
```bash
curl http://localhost:8080/api/dashboard/summary
```
