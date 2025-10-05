# 🐱 Spy Cat API

Golang + Gin + PostgreSQL — fully Dockerized.

## Quick Start


# 1️ Create env and start containers
cp .env.example .env && docker compose up -d --build

# 2️ Apply migrations
docker compose run --rm migrate up

App runs at [http://localhost:8080](http://localhost:8080)

## Swagger UI

Open in browser:
[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)


## ⚙️ Stack

* **Language:** Go 1.24
* **Framework:** Gin
* **Database:** PostgreSQL 16
* **Migrations:** migrate/migrate
* **Docs:** Swagger (swaggo)

## 🗂 Project structure

cmd/                     # Entrypoints
  ├── app/               # Main application
  └── migrator/          # DB migration

docs/                    # Swagger documentation
  ├── docs.go
  ├── swagger.json
  └── swagger.yaml

internal/                # Internal application modules
  ├── app/               # App initialization & router
  ├── config/            # Configuration
  ├── controllers/       # HTTP
  │   ├── http/
  │   │   ├── dto/       # Request/response DTOs
  │   │   └── handlers/  # Route handlers
  ├── helpers/           # Error helpers
  ├── middleware/        # Middleware
  ├── validator/         # Validation
  ├── domain/            # Core domain models 
  ├── lib/               # Shared libraries (PostgreSQL connection, etc.)
  ├── repository/        # Data access  (PostgreSQL repos)
  ├── service/           # Business logic
  └── services_errors/   # Centralized service errors

migrations/              # SQL migrations