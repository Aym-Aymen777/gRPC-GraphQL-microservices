# Elasticsearch Guide (Go Repository – Catalog Service)

## 📌 Overview

This document summarizes how Elasticsearch is used in this repository to implement **CRUD + Search operations** for products.

Goal:
👉 Help you quickly recall **concepts + patterns + best practices**

---

# 🧠 Core Concepts

## 1) Index vs Document

* **Index** → like a database table (`products`)
* **Document** → like a row (JSON object)

```json
{
  "name": "iPhone",
  "price": 1200
}
```

---

## 2) `_id` vs `_source`

* `_id` → document identifier (NOT inside `_source`)
* `_source` → actual stored data

---

## 3) Elasticsearch is NOT your main DB

* MongoDB → source of truth
* Elasticsearch → search engine (read-optimized)

---

# 🏗 Repository Structure

```go
type ElasticSearchRepository struct {
    db *elasticsearch.Client
}
```

* Uses official Go client
* All operations go through `r.db`

---

# 🔌 Initialization

```go
elasticsearch.NewClient(elasticsearch.Config{
    Addresses: []string{url},
})
```

✔ Connects to ES node
✔ Usually injected into service layer

---

# ✍️ CREATE (Index Document)

## Method

```go
CreateProduct(ctx, product)
```

## What happens

* Convert struct → JSON
* Store using `Index()`

## Key points

* `WithDocumentID(product.ID)` → custom ID
* `refresh=true` → immediate visibility (dev only)

---

# 🔍 READ (Single Document)

## Method

```go
GetProductByID(ctx, id)
```

## Flow

1. `Get(index, id)`
2. Decode JSON response
3. Extract `_source`
4. Map to struct

## Important

* `_source` contains data
* `_id` must be manually assigned

---

# 🔎 READ (Multiple Documents)

## Method

```go
GetProductsByIDs(ctx, ids)
```

## Uses

👉 `Mget` (multi-get API)

## Flow

* Send:

```json
{ "ids": ["1", "2", "3"] }
```

* Handle:

  * `found: false` → skip
  * `_source` → actual data

## Key concept

✔ Efficient batch retrieval
✔ Avoids multiple network calls

---

# 🔄 UPDATE (Partial Update)

## Method

```go
UpdateProduct(ctx, product)
```

## Flow

1. Check existence (`Exists`)
2. Send partial update:

```json
{
  "doc": {
    "name": "...",
    "price": ...
  }
}
```

## Why partial update?

* Prevent overwriting full document
* Safer + faster

---

# ❌ DELETE

## Method

```go
DeleteProduct(ctx, id)
```

## Behavior

* Deletes document by `_id`
* If not found → returns `nil` (idempotent)

## Concept

👉 Delete should be safe to repeat

---

# 📄 LIST (Basic Search)

## Method

```go
ListProducts(ctx)
```

## Uses

👉 `Search()` with no query

## Returns

* First N documents (`size=100`)

## Important

* Not scalable for large datasets
* Should use pagination later

---

# 🔍 SEARCH (Full-text)

## Method

```go
SearchForProducts(ctx, query)
```

## Query Used

```json
{
  "multi_match": {
    "query": "iphone",
    "fields": ["name^2", "description"]
  }
}
```

## Key Concepts

### 🔹 Multi-match

Search across multiple fields

### 🔹 Boosting

```json
"name^2"
```

👉 Name is more important than description

---

# 🔁 RESPONSE STRUCTURE

## Standard ES response:

```json
{
  "hits": {
    "hits": [
      {
        "_id": "1",
        "_source": { ... }
      }
    ]
  }
}
```

## Your struct:

```go
type esSearchResponse struct {
    Hits struct {
        Hits []struct {
            ID     string
            Source Product
        }
    }
}
```

✔ Clean parsing
✔ No unsafe casting

---

# ⚠️ COMMON PITFALLS

## ❌ Using `_source.id`

* ID is in `_id`, not `_source`

## ❌ Unsafe type assertions

```go
source["price"].(float64)
```

## ❌ Using ES as primary DB

## ❌ `refresh=true` in production

* Slows performance

---

# ⚡ PERFORMANCE RULES

## ✅ Use Mget for batch reads

## ✅ Avoid refresh in production

## ✅ Use pagination:

```go
WithFrom(offset)
WithSize(limit)
```

---

# 🧠 ARCHITECTURE FLOW

## Correct Data Flow

```text
Create:
MongoDB → Elasticsearch

Update:
MongoDB → Elasticsearch

Delete:
MongoDB → Elasticsearch
```

👉 ES is always **secondary**

---

# 🚀 NEXT LEVEL (TO LEARN)

## 1) Pagination

```go
WithFrom(0)
WithSize(20)
```

## 2) Filters

```json
"range": {
  "price": { "gte": 100, "lte": 1000 }
}
```

## 3) Sorting

```go
WithSort("price:asc")
```

## 4) Autocomplete

* edge_ngram
* custom analyzers

## 5) Bulk indexing

👉 Critical for scaling

---

# 🧾 SUMMARY

| Operation | ES API         | Purpose         |
| --------- | -------------- | --------------- |
| Create    | Index          | Insert document |
| Get       | Get            | Single fetch    |
| Multi Get | Mget           | Batch fetch     |
| Update    | Update         | Partial update  |
| Delete    | Delete         | Remove doc      |
| List      | Search         | Basic fetch     |
| Search    | Search + Query | Full-text       |

---

# 🎯 FINAL MENTAL MODEL

Think of Elasticsearch as:

> ⚡ “A high-speed search layer on top of your database”

NOT:

> ❌ “A replacement for MongoDB”

---

