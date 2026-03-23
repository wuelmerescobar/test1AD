# Library Management System (LMS) – Admin Backend & Frontend

## Project Overview

This project is a **Library Management System (LMS)** built as part of a backend development assessment.
It provides an API and a simple admin frontend to manage:

* Branches
* Books
* Members
* Book availability per branch

The system follows a **clean architecture** using Go and demonstrates production-level features such as middleware, rate limiting, logging, and metrics.

---

## Objective

The goal of this project is to design and implement a backend system that:

* Handles real-world API requests
* Uses middleware for security and performance
* Provides monitoring through metrics
* Demonstrates clean and maintainable code structure

---

## Tech Stack

### Backend

* Go (Golang)
* PostgreSQL
* net/http (standard library)

### Frontend

* HTML
* CSS
* JavaScript (Vanilla)

### Tools

* curl (testing)
* jq (JSON formatting)
* Makefile (automation)

---

## Project Structure

```
test1/
├── cmd/api/main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── router/
│   ├── middleware/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   ├── models/
│   └── validator/
├── migrations/
├── frontend/
├── Makefile
└── README.md
```

---

## Architecture

The system follows a layered structure:

* **Handler Layer** → Handles HTTP requests/responses
* **Service Layer** → Business logic
* **Repository Layer** → Database access
* **Middleware Layer** → Cross-cutting concerns

This ensures:

* Separation of concerns
* Easier testing
* Scalability

---

## Middleware Implemented

The system includes the following middleware:

### 1. Rate Limiting

* Limits number of requests per second
* Returns:

  * `429 Too Many Requests`
  * `Retry-After` header

### 2. Logging

* Structured JSON logs
* Includes method, path, status, latency

### 3. Metrics

* Tracks:

  * Total requests
  * Requests per route
  * Requests per method
  * Error count
  * Average latency
  * Latency per route

Endpoint:

```
GET /metrics
```

### 4. CORS

* Allows frontend (localhost:5500) to access API
* Handles preflight requests (OPTIONS)

### 5. Gzip Compression

* Compresses responses if client supports it

### 6. Recovery Middleware

* Prevents server crash from panics

---

## API Endpoints

### Branches

```
GET    /branches
POST   /branches
```

### Books

```
GET    /books
POST   /books
PUT    /books/{id}
DELETE /books/{id}
```

### Book Availability

```
GET /branches/{id}/books
```

### Members

```
POST   /members
DELETE /members/{id}
GET    /branches/{id}/members
```

### System

```
GET /health
GET /metrics
```

---

## Database

The system uses PostgreSQL with migrations.

Tables include:

* branches
* books
* book_copies
* members

Sample data is inserted automatically.

---

## How to Run the Project

### 1. Set Environment Variables

```bash
export DB_DSN="postgres://postgres:test1Access@localhost:5432/library-test1?sslmode=disable"
```

---

### 2. Setup Database

```bash
make setup
```

---

### 3. Run Server

```bash
go run ./cmd/api/main.go
```

Server runs on:

```
http://localhost:8080
```

---

### 4. Run Frontend

Open:

```
frontend/index.html
```

or use Live Server:

```
http://localhost:5500
```

---

## Testing Features

### Rate Limiting

```
for i in {1..10}; do curl -i http://localhost:8080/branches; done
```

---

### Metrics

```
curl -s http://localhost:8080/metrics | jq
```

---

### Gzip Compression

```
curl -i -H "Accept-Encoding: gzip" http://localhost:8080/branches
```

---

### CORS Preflight

```
curl -i -X OPTIONS http://localhost:8080/branches \
  -H "Origin: http://localhost:5500" \
  -H "Access-Control-Request-Method: GET"
```

---

## Frontend Features

The admin interface allows:

* View branches (button-based UI)
* View books (with pagination)
* Add books to specific branches
* Edit books
* Delete books
* Add members
* Delete members
* View books by branch
* View members by branch

The UI is structured like an **admin dashboard**.

---

## Key Features Implemented

* Clean layered architecture
* Middleware pipeline
* Metrics collection with latency tracking
* Rate limiting for API protection
* JSON structured logging
* Graceful shutdown handling
* Frontend integration

---

## ERD

[text](../../..)

---

## Demo Summary

The system demonstrates:

* API functionality
* Middleware in action
* Frontend interacting with backend
* Real-time metrics tracking

---

## Conclusion

This project demonstrates how to build a **production-ready backend system** with:

* Proper architecture
* Middleware handling
* Monitoring and observability
* API + frontend integration

It provides a solid foundation for expanding into a full-scale application with authentication and role-based access control.

---
