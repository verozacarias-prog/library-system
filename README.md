# Library Management System

Distributed library management system built with NestJS and Go microservices.

## Architecture

Two independent services communicating via HTTP:

```
Client
  │
  ▼ (all requests, including loans)
library-service (NestJS, port 3000)
  │  ├─ manages books, users, auth
  │  └─────────── HTTP ──────────▶ loans-service (Go, port 8081)
  │                                      │  └─ manages loan records
  ▼                                      ▼
postgres-library                   postgres-loans
```

**library-service** (Port 3000) — single entry point for all clients. Manages books, users, authentication, and proxies loan operations to loans-service.

**loans-service** (Port 8081) — internal service, not directly exposed to clients. Manages loan records and validates book availability by calling library-service.

Each service has its own PostgreSQL database — separation of concerns at the data level.

Inter-service communication is authenticated: loans-service generates a short-lived JWT (signed with the shared `JWT_SECRET`, `role: "service"`, 1-minute expiry) on each outbound request to library-service. library-service validates it with the standard `JwtAuthGuard`.

---

## Setup

### Requirements

- Docker Desktop
- Go 1.26+
- Node.js 24+

### Run with Docker Compose

```bash
docker compose up --build
```

This starts 4 containers: library-service, loans-service, postgres-library, postgres-loans.

The loans-service waits for postgres-loans to be healthy before starting (via Docker healthcheck).

The library-service waits for postgres-library to be healthy before starting (via Docker healthcheck).

The loans schema is applied automatically on first run from `loans-service/db/schema.sql`.

### Environment variables

Copy `.env.example` to `.env` and fill in the values before running locally without Docker.

```env
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/library_db
LOANS_SERVICE_URL=http://localhost:8081
JWT_SECRET=your_secret_here
LIBRARY_SERVICE_URL=http://localhost:3000
PORT=3000
```

> **Note:** JWT_SECRET must never be hardcoded in production. Use a secrets manager or external environment variable injection. Both services must share the same JWT_SECRET value.

---

## API Reference

All requests go through **library-service** on port 3000. loans-service (port 8081) is internal only.

### Auth

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/auth/login` | None | Login, returns JWT |

### Users

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/users` | None | Register a new user |
| GET | `/users` | JWT + admin | List all users |
| GET | `/users/:id` | JWT | Get user by ID |
| PATCH | `/users/:id` | JWT | Update user (only admin can change role) |
| DELETE | `/users/:id` | JWT + admin | Delete user |

### Books

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/books` | None | List books with filters and pagination |
| GET | `/books/:id` | None | Get book by ID |
| POST | `/books` | JWT + admin | Create book |
| PATCH | `/books/:id` | JWT + admin | Update book |
| DELETE | `/books/:id` | JWT + admin | Delete book |
| PATCH | `/books/:id/copies` | JWT (internal) | Update available copies — called by loans-service |

**Filters for `GET /books`:** `?author=&genre=&available=true&page=1&limit=10`

### Loans

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/loans` | JWT | Create loan (proxied to loans-service) |
| PATCH | `/loans/:id` | JWT | Return a book (proxied to loans-service) |
| GET | `/loans/users/:userId` | JWT | Active loans for a user |
| GET | `/loans/users/:userId/history` | JWT | Loan history for a user |

### Health

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/health` | None | loans-service health check (port 8081 direct) |

---

## End-to-End Testing

All commands below run against `http://localhost:3000` after `docker compose up --build`.

Replace `$TOKEN` with the value returned by the login step, or use the `export TOKEN=` command shown in Step 3.

### Step 1 — Register an admin user

```bash
curl -s -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Verónica", "email": "vero@test.com", "password": "secret123", "role": "admin"}' \
  | jq
```

Expected: `201` with user object (no password field).

```json
<!-- EVIDENCE: paste response here -->
```

### Step 2 — Register a regular user

```bash
curl -s -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Lector", "email": "lector@test.com", "password": "secret123", "role": "user"}' \
  | jq
```

Expected: `201` with user object.

```json
<!-- EVIDENCE: paste response here -->
```

### Step 3 — Login as admin

```bash
curl -s -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "vero@test.com", "password": "secret123"}' \
  | jq
```

Expected: `200` with `access_token`.

```json
<!-- EVIDENCE: paste response here -->
```

Export the token for the rest of the session:

```bash
export TOKEN="eyJhbGci..."
```

### Step 4 — Create a book (admin only)

```bash
curl -s -X POST http://localhost:3000/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "Clean Code", "author": "Robert Martin", "isbn": "9780132350884", "year": 2008, "genre": "tech", "available_copies": 3}' \
  | jq
```

Expected: `201` with book object. Note the `id` for the next steps.

```json
<!-- EVIDENCE: paste response here -->
```

### Step 5 — List books with filters

```bash
curl -s "http://localhost:3000/books?author=Martin&genre=tech&available=true&page=1&limit=10" | jq
```

Expected: `200` with `{ data: [...], total: 1 }`.

```json
<!-- EVIDENCE: paste response here -->
```

### Step 6 — Create a loan

```bash
curl -s -X POST http://localhost:3000/loans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"user_id": 1, "book_id": 1}' \
  | jq
```

Expected: `201` with loan object (`status: "active"`). loans-service validated book availability with library-service before persisting.

```json
<!-- EVIDENCE: paste response here -->
```

### Step 7 — Verify available copies decremented

```bash
curl -s http://localhost:3000/books/1 | jq '.available_copies'
```

Expected: `2` (was 3).

```
<!-- EVIDENCE: paste output here -->
```

### Step 8 — Get active loans for user

```bash
curl -s http://localhost:3000/loans/users/1 \
  -H "Authorization: Bearer $TOKEN" \
  | jq
```

Expected: `200` with array containing the active loan.

```json
<!-- EVIDENCE: paste response here -->
```

### Step 9 — Return the book

```bash
curl -s -X PATCH http://localhost:3000/loans/1 \
  -H "Authorization: Bearer $TOKEN" \
  | jq
```

Expected: `200` with loan object (`status: "returned"`). Available copies incremented back to 3.

```json
<!-- EVIDENCE: paste response here -->
```

### Step 10 — Verify copies restored

```bash
curl -s http://localhost:3000/books/1 | jq '.available_copies'
```

Expected: `3`.

```
<!-- EVIDENCE: paste output here -->
```

### Step 11 — View loan history

```bash
curl -s http://localhost:3000/loans/users/1/history \
  -H "Authorization: Bearer $TOKEN" \
  | jq
```

Expected: `200` with array containing the returned loan.

```json
<!-- EVIDENCE: paste response here -->
```

### Step 12 — Error cases

**Borrow a book with no copies available:**

Create 2 more loans first to exhaust the remaining copies:

```bash
curl -s -X POST http://localhost:3000/loans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"user_id": 2, "book_id": 1}' | jq

curl -s -X POST http://localhost:3000/loans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"user_id": 3, "book_id": 1}' | jq
```

Then try one more:

```bash
curl -s -X POST http://localhost:3000/loans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"user_id": 4, "book_id": 1}' \
  | jq
```

Expected: `409 Conflict`.

```json
<!-- EVIDENCE: paste response here -->
```

**Try to create a book without admin role:**

```bash
export LECTOR_TOKEN=$(curl -s -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "lector@test.com", "password": "secret123"}' | jq -r '.access_token')

curl -s -X POST http://localhost:3000/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $LECTOR_TOKEN" \
  -d '{"title": "Test", "author": "X", "isbn": "1234567890123", "year": 2024, "genre": "test", "available_copies": 1}' \
  | jq
```

Expected: `403 Forbidden`.

```json
<!-- EVIDENCE: paste response here -->
```

**Access a protected endpoint without token:**

```bash
curl -s -X POST http://localhost:3000/loans \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "book_id": 1}' \
  | jq
```

Expected: `401 Unauthorized`.

```json
<!-- EVIDENCE: paste response here -->
```

---

## Technical Decisions

### loans-service (Go)

**HTTP Framework: Chi over Gin**
Chi is closer to the standard library, more idiomatic Go, and easier to test. Gin adds magic that hides what the router is doing. Since test coverage is a goal, Chi was the better fit.

**Database driver: pgx over GORM**
GORM abstracts SQL in ways that can hide performance problems and errors. Using `pgx` directly with `database/sql` shows explicit control over queries and transactions. In a technical assessment, this communicates more about the candidate than GORM would.

**Two separate PostgreSQL instances**
Each service owns its data completely. A shared database creates coupling at the data layer — if library-service changes its schema, loans-service breaks. Two instances enforce true independence. Trade-off: higher resource usage. Justified for a microservices architecture.

**Multi-stage Docker build**
The final image contains only the compiled binary and Alpine Linux (~5MB + binary). No Go compiler, no source code. Smaller attack surface, faster deploy.

**Constants centralized per layer**
SQL queries live in `repository/constants.go`. Domain status values (`active`, `returned`) and copy actions (`increment`, `decrement`) live in `internal/constants.go`. Error variables live in `errors.go` files per layer (repository, service, handler, clients).

**BookServiceClient as interface defined in service layer**
Following dependency inversion: the service layer defines the interface it needs (`BookServiceClient`), and the concrete implementation (`clients/library_client.go`) implements it implicitly (Go duck typing). This keeps the service layer independent of HTTP details and makes it fully testable with mocks.

**Factory pattern (internal/app.go)**
All infrastructure dependencies (pool, repository, service) are constructed in one place. `main.go` only creates the HTTP layer (handler, router) and wires it together. The factory does not include the HTTP framework to avoid coupling infrastructure with transport.

**Schema with partial unique index**
```sql
CREATE UNIQUE INDEX ON loans(user_id, book_id) WHERE status = 'active';
```
A user cannot have the same book twice on active loan simultaneously. Once returned, they can borrow it again. A simple UNIQUE constraint would prevent that.

**PATCH /loans/:id for returns**
Instead of a `/loans/:id/return` verb-in-URL pattern, a `PATCH /loans/:id` was used. The service layer method `BookReturned` always sets the status to `returned` — the status is a consequence of the operation, not an input from the client.

**Inter-service authentication with short-lived JWT**
loans-service generates a JWT signed with the shared `JWT_SECRET` (`role: "service"`, 1-minute expiry) on each outbound call to library-service. library-service validates it with the standard `JwtAuthGuard` — no new authentication mechanism needed. The short expiry limits the blast radius of a leaked token in transit.

---

### library-service (NestJS)

**TypeORM over Prisma**
TypeORM is more mature for complex relational mappings and has better TypeScript integration with NestJS decorators. Both were valid choices — TypeORM was selected for stability.

**autoLoadEntities: true**
Instead of listing every entity in `TypeOrmModule.forRoot()`, each module registers its own entity with `forFeature()`. Scales better as the number of entities grows.

**synchronize: true**
TypeORM automatically creates/updates tables based on entity definitions. Appropriate for development and this assessment. In production, use migrations (`typeorm migration:generate` + `typeorm migration:run`).

**JWT with Passport**
`@nestjs/passport` + `passport-jwt` is the standard NestJS authentication pattern. The `JwtStrategy` validates the token from the `Authorization: Bearer` header. `JwtAuthGuard` is applied per-endpoint with `@UseGuards()`.

**Role-based access control with custom decorator + guard**
`@Roles('admin')` decorator stores metadata on the endpoint. `RolesGuard` reads that metadata via `Reflector` and compares it against the role in the JWT payload. Endpoints without `@Roles` are accessible to any authenticated user. Only admins can change a user's role via `PATCH /users/:id`.

**Input validation with class-validator**
All endpoints use typed DTOs with `class-validator` decorators. A global `ValidationPipe` (whitelist + forbidNonWhitelisted) is applied in `main.ts`, rejecting unknown fields and validating types at the boundary.

**Passwords hashed with bcrypt (salt rounds: 10)**
Passwords are never stored or returned in plain text. The `Omit<User, 'password'>` TypeScript type ensures the password field is excluded at the type level from all service responses.

**Pagination with QueryBuilder**
`findAll` uses TypeORM's `createQueryBuilder` to support dynamic filters (author, genre, availability) combined with `skip/take` for pagination. Returns `{ data, total }` so the client can implement pagination UI.

---

## Testing

```bash
# loans-service
cd loans-service
go test ./internal/service/...
go test ./internal/handler/...

# library-service
cd library-service
npm test
```

**loans-service:** 6 unit tests — service layer (CreateLoan happy path, book not available, BookReturned, GetActiveLoans) and handler layer (CreateLoan 201, invalid body 400). Mocks used for repository and LibraryClient — no database required.

**library-service:** 6 unit tests — BooksService (create, findOne not found, updateCopies) and AuthService (valid login, user not found, wrong password). Mocks used for TypeORM repository and JwtService.

---

## What was intentionally left out

- **gRPC between services:** HTTP was kept for simplicity and consistency. gRPC would be the natural next step to reduce inter-service latency.
- **TypeORM migrations:** `synchronize: true` is used for this assessment. In production, replace with `typeorm migration:generate` + `typeorm migration:run` to get versioned, reversible schema changes.
- **User existence validation in loans-service:** loans-service does not call library-service to verify that the userId is valid before creating a loan. A dedicated `/users/:id` validation call could be added, but adds latency for the common path.
- **Rate limiting:** out of scope for this assessment.
- **Frontend:** explicitly excluded by the assessment.
- **Kubernetes / service mesh:** explicitly excluded by the assessment.
- **Repository integration tests:** require a running database instance. Prioritized unit tests for business logic coverage.

---

## Project Structure

```
library-system/
├── docker-compose.yml
├── .env.example
├── library-service/              # NestJS — books, users, auth, loans proxy
│   ├── src/
│   │   ├── books/
│   │   │   ├── book.entity.ts
│   │   │   ├── books.module.ts
│   │   │   ├── books.service.ts
│   │   │   ├── books.service.spec.ts
│   │   │   ├── books.controller.ts
│   │   │   └── dto/
│   │   │       ├── create-book.dto.ts
│   │   │       └── update-book.dto.ts
│   │   ├── users/
│   │   │   ├── user.entity.ts
│   │   │   ├── users.module.ts
│   │   │   ├── users.service.ts
│   │   │   ├── users.controller.ts
│   │   │   └── dto/
│   │   │       ├── create-user.dto.ts
│   │   │       └── update-user.dto.ts
│   │   ├── loans/                    # proxy to loans-service
│   │   │   ├── loans.module.ts
│   │   │   └── loans.controller.ts
│   │   └── auth/
│   │       ├── auth.module.ts
│   │       ├── auth.service.ts
│   │       ├── auth.service.spec.ts
│   │       ├── auth.controller.ts
│   │       ├── jwt.strategy.ts
│   │       ├── jwt-auth.guard.ts
│   │       ├── roles.guard.ts
│   │       └── roles.decorator.ts
│   └── Dockerfile
└── loans-service/                # Go — loans management
    ├── cmd/api/main.go
    ├── internal/
    │   ├── app.go
    │   ├── constants.go
    │   ├── clients/
    │   │   ├── library_client.go
    │   │   └── errors.go
    │   ├── handler/
    │   │   ├── loan_handler.go
    │   │   ├── loan_handler_test.go
    │   │   ├── mocks_test.go
    │   │   ├── routes.go
    │   │   └── errors.go
    │   ├── service/
    │   │   ├── loan_service.go
    │   │   ├── loan_service_test.go
    │   │   ├── mocks_test.go
    │   │   ├── library_client.go
    │   │   └── errors.go
    │   ├── repository/
    │   │   ├── loan_repository.go
    │   │   ├── db.go
    │   │   ├── constants.go
    │   │   └── errors.go
    │   └── model/
    │       ├── loan.go
    │       └── book.go
    ├── db/schema.sql
    ├── Dockerfile
    └── go.mod
```
