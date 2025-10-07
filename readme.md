## 🐱Spy Cat API

A REST API in Go (Gin) with PostgreSQL.
The project ships with Docker Compose (app + Postgres + migrations),
clean module structure, database migrations, and Swagger/OpenAPI docs.

## Quick Start
Everything can be set up with just two commands.

``` text
1. Copy environment and start containers
cp .env.example .env && docker compose up -d --build

2. Apply database migrations
docker compose run --rm migrate up

[App runs at](http://localhost:8080)
[Swagger UI](http://localhost:8080/swagger/index.html)
```
## Full Setup (for any OS)
```text
git clone https://github.com/<your-username>/spy-cat.git
cd spy-cat

# 1. Copy environment file
cp .env.example .env   # Linux / macOS
# (Windows PowerShell) Copy-Item .env.example .env

# 2. Build and start containers
docker compose up -d --build

# 3. Apply migrations
docker compose run --rm migrate up

# 4. (Optional) Seed demo data
docker compose run --rm seed
```
## 🐾 Seeding demo data
``` text
You can optionally fill the database with demo cats from
**TheCatAPI** (https://api.thecatapi.com/v1/breeds)

1. Make sure services are running
docker compose up -d --build
2. Run the seeder
docker compose run --rm seed

Example output:
200! Cornish Rex (exp=15, salary=2139)
200! Cymric (exp=14, salary=4386)
Done.
```

## ⚙️ Stack
* **Language:** Go 1.24
* **Framework:** Gin
* **Database:** PostgreSQL 16
* **Migrations:** migrate/migrate
* **Docs:** Swagger (swaggo)

## 📚 API Documentation

```text
The API uses auto-generated OpenAPI (Swagger) schema.
You can test endpoints directly in the browser:

Swagger : http://localhost:8080/swagger/index.html
Raw spec: http://localhost:8080/swagger/doc.json
```

## Environment
Use .env.example as a base:

```text
APP_PORT=8080

DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=changeme
DB_NAME=spy_cat
DB_SSLMODE=disable
```
## 🗂 Project structure

``` text
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
```

## Development Tips
``` text
1) View logs
docker compose logs -f app

2) Apply / rollback migrations
docker compose run --rm migrate up
docker compose run --rm migrate down 1

3) Rebuild containers
docker compose up -d --build

4) Stop and clean up
docker compose down
```