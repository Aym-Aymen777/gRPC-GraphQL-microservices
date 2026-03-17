# gRPC + GraphQL Microservices

A learning project demonstrating how to build a **microservices architecture using Go, gRPC, GraphQL, and MySQL**.

The system is composed of independent services communicating through **gRPC**, while a **GraphQL gateway** exposes a unified API for clients.

---

# Architecture Overview

```
Client (Web / Mobile)
        ↓
     GraphQL API
        ↓
     gRPC Clients
        ↓
 ┌───────────────┬───────────────┬───────────────┐
 │ Account       │ Catalog       │ Order         │
 │ Microservice  │ Microservice  │ Microservice  │
 └───────────────┴───────────────┴───────────────┘
        ↓               ↓               ↓
      MySQL        ElasticSearch       MySQL
```

Each service is **independent and responsible for its own data**.

---

# Services

## Account Service

Manages user accounts.

Responsibilities:

* Create accounts
* Fetch account information
* List accounts

Documentation:

```
/account/tutorial.md
```

---

## Catalog Service

Handles products or items available in the system.

Responsibilities:

* Create catalog items
* Retrieve items
* List items

Documentation:

```
/catalog/tutorial.md
```

---

## Order Service

Manages orders and transactions.

Responsibilities:

* Create orders
* Retrieve order history
* Manage order status

Documentation:

```
/order/tutorial.md
```

---

## GraphQL Gateway

Provides a **single API endpoint for clients**.

Responsibilities:

* Translate GraphQL queries into gRPC calls
* Aggregate responses from services
* Simplify client communication

Documentation:

```
/graphql/tutorial.md
```

---

# Service Architecture

Each microservice follows the same internal structure.

```
Client
   ↓
gRPC Server
   ↓
Service (Business Logic)
   ↓
Repository (Database Access)
   ↓
Database
```

Layers:

| Layer      | Responsibility            |
| ---------- | ------------------------- |
| Server     | Exposes gRPC endpoints    |
| Service    | Contains business logic   |
| Repository | Executes database queries |
| Database   | Stores service data       |

---

# Project Structure

```
gRPC-GraphQL-microservices
│
├── account
│   ├── repository
│   ├── service
│   ├── server
│   ├── types
│   ├── proto
│   ├── tutorial.md
│   └── main.go
│
├── catalog
│   ├── repository
│   ├── service
│   ├── server
│   ├── types
│   ├── proto
│   ├── tutorial.md
│   └── main.go
│
├── order
│   ├── repository
│   ├── service
│   ├── server
│   ├── types
│   ├── proto
│   ├── tutorial.md
│   └── main.go
│
├── graphql
│   ├── resolvers
│   ├── schema
│   ├── clients
│   ├── tutorial.md
│   └── main.go
│
└── README.md
```

---

# Technologies

Main technologies used in this project:

* **Go**
* **gRPC**
* **GraphQL**
* **MySQL**
* **Protocol Buffers**

---

# Communication Flow

Example flow for creating an account.

```
GraphQL Mutation
        ↓
GraphQL Resolver
        ↓
gRPC Client
        ↓
Account Service
        ↓
Service Layer
        ↓
Repository
        ↓
MySQL
```

---

# Running the Services

Each service runs independently.

Example:

```
go run account/main.go
go run catalog/main.go
go run order/main.go
go run graphql/main.go
```

Services communicate through **gRPC endpoints**.

---

# Learning Goals

This project demonstrates:

* Microservices architecture
* Service-to-service communication with gRPC
* API gateway using GraphQL
* Repository pattern
* Clean service layering

---

# Documentation

Detailed tutorials for each service:

```
account/tutorial.md
catalog/tutorial.md
order/tutorial.md
graphql/tutorial.md
```

---

# Future Improvements

Possible improvements for the project:

* Dockerization
* Service discovery
* Authentication
* Distributed tracing
* Message queues

---

# License

This project is for **educational purposes**.
