# 🚦 Distributed Rate Limiter in Go

> A production-inspired distributed rate limiter built with Go and Redis, designed to demonstrate backend engineering principles, distributed system concepts, and production-ready software design.

---

## 📖 Overview

Modern distributed applications often run multiple instances of the same service behind a load balancer. While this improves scalability and availability, it also introduces a challenge:

> **How do multiple servers enforce the same request limits for a user?**

A simple in-memory rate limiter works only on a single server. Once multiple instances are introduced, each server maintains its own counters, leading to inconsistent rate limiting.

This project solves that problem by implementing a **distributed rate limiter** where multiple application instances share rate-limiting state through **Redis**.

The goal is not only to build a working service, but also to understand the design decisions, trade-offs, and engineering practices involved in building production systems.

---

# 🎯 Project Goals

This project aims to demonstrate:

* Building production-inspired backend services in Go
* Designing scalable distributed systems
* Implementing efficient rate limiting algorithms
* Writing clean, maintainable Go code
* Following idiomatic Go project structure
* Using Redis as distributed shared state
* Containerizing applications with Docker
* Writing comprehensive tests
* Documenting architectural decisions

---

# ✨ Features

### Core Features

* Distributed rate limiting
* Redis-backed shared state
* HTTP REST API
* Configurable limits
* Docker & Docker Compose support
* Graceful shutdown
* Health endpoints
* Structured configuration
* Production-style project layout

---

### Engineering Features

* Clean Architecture
* Dependency Injection
* Interface-based design
* Configuration management
* Error handling
* Concurrency-safe implementation
* Extensible architecture
* Unit & Integration Tests

---

# 🏗 Architecture

```
                    Client
                       │
                       ▼
              +----------------+
              | Load Balancer  |
              +----------------+
                 │    │    │
        ┌────────┘    │    └────────┐
        ▼             ▼             ▼
+---------------+ +---------------+ +---------------+
| Go Instance A | | Go Instance B | | Go Instance C |
+---------------+ +---------------+ +---------------+
         │              │                 │
         └──────────────┼─────────────────┘
                        ▼
                +----------------+
                |     Redis      |
                | Shared State   |
                +----------------+
```

Each server instance shares rate-limiting state through Redis, ensuring consistent enforcement regardless of which server handles a request.

---

# 📂 Project Structure

```
distributed-rate-limiter/

├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── api/
│   ├── config/
│   ├── limiter/
│   ├── middleware/
│   ├── storage/
│   └── redis/
│
├── tests/
│
├── docs/
│
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
├── .env.example
└── README.md
```

---

# 🛠 Tech Stack

| Category             | Technology            |
| -------------------- | --------------------- |
| Language             | Go                    |
| Cache / Shared State | Redis                 |
| HTTP                 | net/http              |
| Containerization     | Docker                |
| Orchestration        | Docker Compose        |
| Configuration        | Environment Variables |
| Testing              | Go Testing Package    |

---

# 📚 Concepts Covered

This project explores several backend engineering topics:

* Distributed Systems
* Rate Limiting
* Redis
* HTTP Middleware
* Interfaces
* Dependency Injection
* Concurrency
* Atomic Operations
* Docker
* Configuration Management
* Graceful Shutdown
* Error Handling
* Testing
* Clean Architecture

---

# 🚀 Rate Limiting Algorithm

> **Initial Implementation:** Token Bucket

The Token Bucket algorithm allows short bursts of traffic while maintaining a controlled average request rate.

Future versions may include:

* Fixed Window
* Sliding Window Counter
* Sliding Window Log
* Leaky Bucket

along with a comparison of their trade-offs.

---

# ⚙ Configuration

Configuration is managed using environment variables.

Example:

```env
PORT=8080

REDIS_ADDR=localhost:6379

DEFAULT_RATE_LIMIT=10

RATE_LIMIT_WINDOW=1m
```

---

# 🚀 Running Locally

## Prerequisites

* Go 1.24+
* Docker
* Docker Compose

---

## Clone

```bash
git clone https://github.com/<your-username>/distributed-rate-limiter.git

cd distributed-rate-limiter
```

---

## Install Dependencies

```bash
go mod tidy
```

---

## Start Redis

```bash
docker compose up redis
```

---

## Start Application

```bash
go run ./cmd/server
```

---

## Or Run Everything

```bash
docker compose up --build
```

---

# 🧪 Testing

Run unit tests:

```bash
go test ./...
```

Run with coverage:

```bash
go test ./... -cover
```

---

# 📡 API

## Health Check

```
GET /health
```

Response

```json
{
    "status": "ok"
}
```

---

## Rate Limit Check *(planned)*

```
POST /v1/check
```

Request

```json
{
    "key": "user123"
}
```

Response

```json
{
    "allowed": true,
    "remaining": 8,
    "retry_after": 0
}
```

---

# 🧪 Local Verification Checklist

* [ ] Project builds successfully
* [ ] Redis starts correctly
* [ ] Application starts without errors
* [ ] Health endpoint responds
* [ ] Rate limiting works as expected
* [ ] Multiple requests return HTTP 429 after limit
* [ ] Docker Compose runs successfully
* [ ] Unit tests pass

---

# 📈 Future Improvements

* Lua Script based atomic operations
* Sliding Window implementation
* Multiple rate-limiting algorithms
* Prometheus metrics
* Grafana dashboards
* Request tracing
* Distributed benchmarking
* Kubernetes deployment
* Horizontal scaling examples
* Authentication
* Per-user and per-endpoint limits
* Configurable policies
* Redis Cluster support
* CI/CD with GitHub Actions

---

# 📖 Learning Outcomes

By building this project, I aim to gain a practical understanding of:

* Designing distributed systems
* Building production-grade backend services
* Writing idiomatic Go
* Designing clean APIs
* Working with Redis
* Applying software engineering best practices
* Containerizing applications
* Testing backend services
* Making architectural trade-offs

---

# 🤝 Contributing

Contributions, suggestions, and discussions are welcome.

If you have ideas for improving the architecture, implementation, or documentation, feel free to open an issue or submit a pull request.

---

# 📄 License

This project is licensed under the MIT License.

---

## ⭐ Acknowledgements

This project is inspired by the engineering challenges encountered while building scalable backend systems and payment infrastructure. It is intended as a learning-focused implementation that emphasizes clean design, production-inspired practices, and continuous improvement over feature completeness.
