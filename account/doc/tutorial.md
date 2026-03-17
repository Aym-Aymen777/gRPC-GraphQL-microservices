# Account Microservice – Tutorial

This document explains the **steps used to build the Account microservice** in the `gRPC + GraphQL` architecture.

Flow of the service:

```
GraphQL Mutation/Query
        ↓
      Client
        ↓
    gRPC Server
        ↓
      Service
        ↓
    Repository
        ↓
     Database
```

---

# Step 1 — Create Database Schema

Create the SQL schema for accounts.

```sql
CREATE TABLE accounts (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL
);
```

---

# Step 2 — Create Domain Types

Create a shared domain model.

**types/account.go**

```go
package types

type Account struct {
    ID       string
    Username string
    Email    string
}
```

---

# Step 3 — Create Repository Interface

Define database operations.

**repository/repository.go**

```go
type Repository interface {
    Close()
    PutAccount(ctx context.Context, account *types.Account) error
    GetAccount(ctx context.Context, id string) (*types.Account, error)
    ListAccounts(ctx context.Context, skip, limit uint64) ([]*types.Account, error)
}
```

This interface defines the data layer contract.

---

# Step 4 — Implement MySQL Repository

Create the repository implementation.

**repository/mysql.go**

Key responsibilities:

- Open database connection
- Execute SQL queries
- Map database rows to domain types

Example functions:

```
NewMySQLRepository()
PutAccount()
GetAccount()
ListAccounts()
Close()
```

Repository flow:

```
Service → Repository → MySQL
```

---

# Step 5 — Create Service Layer

The service layer contains **business logic**.

**service/service.go**

Example interface:

```go
type Service interface {
    PostAccount(ctx context.Context, username, email string) (*types.Account, error)
    GetAccount(ctx context.Context, id string) (*types.Account, error)
    ListAccounts(ctx context.Context, skip, limit uint64) ([]*types.Account, error)
}
```

Responsibilities:

- Validate input
- Generate IDs
- Call repository methods

Service flow:

```
gRPC Server → Service → Repository
```

---

# Step 6 — Define gRPC API

Create the gRPC definitions.

**proto/account.proto**

Example:

```proto
service AccountService {
  rpc PostAccount(PostAccountRequest) returns (AccountResponse);
  rpc GetAccount(GetAccountRequest) returns (AccountResponse);
  rpc ListAccounts(ListAccountsRequest) returns (ListAccountsResponse);
}
```

Then generate Go code:

```
protoc --go_out=. --go-grpc_out=. proto/account.proto
```

---

# Step 7 — Implement gRPC Server

Create the server implementation.

**server/grpc.go**

Responsibilities:

- Receive RPC requests
- Call service layer
- Return responses

Flow:

```
Client → gRPC Server → Service → Repository
```

---

# Step 8 — Create Service Entry Point

Create the service bootstrap.

**main.go**

Responsibilities:

- Initialize repository
- Initialize service
- Start gRPC server

Startup flow:

```
Start Service
     ↓
Connect Database
     ↓
Start gRPC Server
```

---

# Project Structure

Example structure of the account service.

```
account
 |_proto
   |_account.proto
 |_repository
 |_mysql.go
 |_service
 |_service.go
 |_server.go
 |_types
   |_account.go
 |_main.go
```

---

# Notes

- Repository handles **database logic**
- Service handles **business logic**
- gRPC server exposes **RPC endpoints**
- GraphQL service will call this service through **gRPC client**
- The difference between **Service** interface and **Repository** is that one handle the datbase logic and the other use them in addition of the bussiness logic



# GPT Note

gRPC Server Blocking Behavior

When running the account microservice, the program may appear to stop at:

```go
 account.ListenGRPC(service, cfg.GRPCPort)```

This is expected behavior.

Inside ListenGRPC, the call:

```go
serv.Serve(listener)
```

starts the gRPC server and blocks the main goroutine while waiting for incoming requests. Since the server must continuously listen for connections, this function never returns during normal operation.

Example implementation:

```go
func ListenGRPC(s Service, port string) error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	log.Println("gRPC server listening on", port)

	server := grpc.NewServer()
	pb.RegisterAccountServiceServer(server, &grpcServer{Service: s})

	return server.Serve(listener)
}
```

Usage in main.go:
```go
log.Printf("Account service is connecting to MySQL ..... ⏳")
service := account.NewAccountService(repo)
log.Printf("Account service is connected to MySQL ✅")
log.Printf("Account service is running on port %s 🎯", cfg.GRPCPort)
account.ListenGRPC(service, cfg.GRPCPort)
```

```go
if err := account.ListenGRPC(service, cfg.GRPCPort); err != nil {
	log.Fatal(err)
}
```

# Important notes

- Serve() is blocking by design.

- The service is running correctly even if the terminal shows no further output.

- Add logging before starting the server to confirm it started.

- The port must be formatted like :50051.

- Example runtime log:

- gRPC server listening on :50051

- The service will then wait indefinitely for gRPC requests.

---
