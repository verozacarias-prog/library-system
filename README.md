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
- Go 1.21+
- Node.js 20+

### Run with Docker Compose

```bash
docker compose up --build
```

This starts 4 containers: library-service, loans-service, postgres-library, postgres-loans.

The loans-service waits for postgres-loans to be healthy before starting (via Docker healthcheck).

The loans schema is applied automatically on first run from `loans-service/db/schema.sql`.

> **Note on data persistence:** Use `docker compose down` (without `-v`) to stop containers while preserving database data. Only use `docker compose down -v` if you want to reset all data. Named volumes are declared in `docker-compose.yml` to ensure persistence across restarts.

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
| GET | `/health` | None | library-service health check (port 3000) |
| GET | `/health` | None | loans-service health check (port 8081 direct) |

---

## Complete Flow Example

Step-by-step curl commands for the full loan lifecycle.

### 1. Login as admin

```bash
TOKEN=$(curl -s -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@library.com","password":"adminpass"}' \
  | jq -r '.access_token')
```

### 2. Check a book (no auth required)

```bash
curl -s http://localhost:3000/books/1 | jq '{id,title,author,available_copies}'
# {"id":1,"title":"The Go Programming Language","author":"Alan Donovan","available_copies":3}
```

### 3. Create a loan

```bash
curl -s -X POST http://localhost:3000/loans \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"book_id":1}' | jq .
# {"id":1,"user_id":1,"book_id":1,"loaned_at":"...","returned_at":null,"status":"active"}
```

### 4. Verify copies decremented

```bash
curl -s http://localhost:3000/books/1 | jq .available_copies
# 2  (was 3)
```

### 5. Return the book

```bash
curl -s -X PATCH http://localhost:3000/loans/1 \
  -H "Authorization: Bearer $TOKEN" | jq .
# {"id":1,"user_id":1,"book_id":1,"loaned_at":"...","returned_at":"...","status":"returned"}
```

### 6. Verify copies restored

```bash
curl -s http://localhost:3000/books/1 | jq .available_copies
# 3  (restored)
```

### 7. View loan history

```bash
curl -s http://localhost:3000/loans/users/1/history \
  -H "Authorization: Bearer $TOKEN" | jq .
# [{"id":1,"status":"returned",...}]
```

---

## End-to-End Testing

All tests run against `http://localhost:3000` after `docker compose up --build`.

### 1. Seed the database

A seed script creates the initial test data: 2 admin users, 1 regular user, 4 books, and 3 loans.

```bash
chmod +x scripts/seed.sh
./scripts/seed.sh
```

Expected output: 3 users created, 4 books created, 3 loans created.

### 2. Verify database state

Confirm data was persisted correctly in both databases:

```bash
# Users and books (library-service DB)
docker exec -it library-system-postgres-library-1 psql -U postgres -d library_db \
  -c "SELECT id, name, email, role FROM users;"

docker exec -it library-system-postgres-library-1 psql -U postgres -d library_db \
  -c "SELECT id, title, author, available_copies FROM books;"

# Loans (loans-service DB)
docker exec -it library-system-postgres-loans-1 psql -U postgres -d loans_db \
  -c "SELECT * FROM loans;"
```

### 3. Run CRUD and validation tests (books & users)

Tests UPDATE and DELETE operations for books and users, covering all validation edge cases (invalid fields, missing auth, wrong role, not found, etc.).

```bash
chmod +x scripts/test_update_delete.sh
./scripts/test_update_delete.sh
```

Expected: each request shows its HTTP status and response body. All status codes match expectations. Captured expected responses are in `scripts/test_results_update_delete.txt`.

### 4. Run full functional tests (auth, loans, error cases)

Tests the complete loan lifecycle: create loan → verify copies decremented → return book → verify copies restored → view history. Also covers error cases: no copies available, duplicate loan, missing auth, wrong role.

```bash
chmod +x scripts/test_full.sh
./scripts/test_full.sh
```

Expected highlights:
- Loan creation returns `201` with `status: "active"`
- Available copies decrement after loan, increment after return
- Second loan for same user + book returns `409 Conflict`
- Exhausting all copies returns `409 Conflict`
- Missing token returns `401 Unauthorized`
- Non-admin creating a book returns `403 Forbidden`

Captured expected responses are in `scripts/test_results_full.txt`.

---

## CI

GitHub Actions runs the unit test suite for both services on every push and pull request.

```
.github/workflows/ci.yml
├── test-library-service   → npm ci && npm test  (Node.js 20)
└── test-loans-service     → go test ./internal/...  (Go 1.21)
```

No database or Docker is required — all tests use mocks.

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

**library-service unavailable → 503**
Network errors from `http.Client.Do()` are wrapped with `ErrLibraryServiceUnavailable` using `fmt.Errorf("%w: %w", ...)`, so `errors.Is()` can unwrap them at the handler layer and return a clean `503 Service Unavailable` instead of a raw error string.

**Distributed transaction handling**
Loan creation and return both involve two systems (loans DB + library-service copy count). To maintain consistency:
- On **loan creation**: copies are decremented first, then the loan is persisted. If the DB insert fails, a compensating increment call is made to restore the count.
- On **loan return**: the loan is fetched first (to get the bookID), copies are incremented, then the loan status is updated. If the status update fails, a compensating decrement call restores the count.

This is best-effort consistency — not a distributed transaction. A failure in the compensating call would leave the systems out of sync, which is a known trade-off documented under "What was intentionally left out."

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
All endpoints use typed DTOs with `class-validator` decorators. A global `ValidationPipe` (whitelist + forbidNonWhitelisted) is applied in `main.ts`, rejecting unknown fields and validating types at the boundary. This includes the loans proxy — `POST /loans` validates `user_id` and `book_id` via `CreateLoanDto` before forwarding to loans-service.

**Atomic copy update**
`updateCopies` uses a single `UPDATE ... SET available_copies = available_copies + $delta WHERE available_copies + $delta >= 0 RETURNING *` query instead of a read-modify-write cycle. This eliminates the race condition where two concurrent loan requests could both read the same copy count and produce an incorrect result.

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

**loans-service:** unit tests — service layer (CreateLoan happy path, book not available, BookReturned, GetActiveLoans) and handler layer (CreateLoan 201, invalid body 400). Mocks used for repository and LibraryClient — no database required.

**library-service:** unit tests — BooksService (create, findOne not found, updateCopies) and AuthService (valid login, user not found, wrong password). Mocks used for TypeORM repository and JwtService.

Tests also run automatically on every push via GitHub Actions (see [CI](#ci) section).

---

## Role permissions

| Action | admin | user |
|--------|-------|------|
| Register account | ✅ | ✅ |
| Login | ✅ | ✅ |
| List books / Get book | ✅ | ✅ |
| Create / Update / Delete book | ✅ | ❌ |
| List all users | ✅ | ❌ |
| Get / Update own user | ✅ | ✅ |
| Change a user's role | ✅ | ❌ |
| Delete user | ✅ | ❌ |
| Create loan | ✅ | ✅ |
| Return book | ✅ | ✅ |
| View active loans for a user | ✅ | ✅ |
| View loan history for a user | ✅ | ✅ |

The assessment specifies "read and their loans" for regular users. This was interpreted as: any authenticated user can manage loans (create and return), since a borrowing system where only admins can register loans would not be functional for end users. Loan endpoints do not restrict by role — only by authentication.

**Suggested improvement:** `GET /loans/users/:userId` currently allows any authenticated user to view another user's loans. A natural next step would be to restrict this so regular users can only query their own loans (validating that `:userId` matches the `sub` from the JWT), while admins retain access to any user's loans.

---

## What was intentionally left out

- **gRPC between services:** HTTP was kept for simplicity and consistency. gRPC would be the natural next step to reduce inter-service latency.
- **TypeORM migrations:** `synchronize: true` is used for this assessment. In production, replace with `typeorm migration:generate` + `typeorm migration:run` to get versioned, reversible schema changes.
- **User existence validation in loans-service:** loans-service does not call library-service to verify that the userId is valid before creating a loan. A dedicated `/users/:id` validation call could be added, but adds latency for the common path.
- **Loan access restriction by ownership:** any authenticated user can currently view any other user's loans. Restricting `GET /loans/users/:userId` to the owner (or admin) is a straightforward improvement not implemented in this version.
- **Compensating transaction reliability:** if the compensating call (copy revert) fails after a DB error, the two systems will be out of sync. A production implementation would use an outbox pattern or saga coordinator. Out of scope for this assessment.
- **Rate limiting:** out of scope for this assessment.
- **Frontend:** explicitly excluded by the assessment.
- **Kubernetes / service mesh:** explicitly excluded by the assessment.
- **Repository integration tests:** require a running database instance. Prioritized unit tests for business logic coverage.

---

## Project Structure

```
library-system/
├── .github/
│   └── workflows/
│       └── ci.yml                # GitHub Actions — runs tests on push/PR
├── docker-compose.yml
├── .env.example
├── scripts/
│   ├── seed.sh                   # creates test users, books, loans
│   ├── test_update_delete.sh     # CRUD + validation tests for books and users
│   └── test_full.sh              # auth, loans lifecycle, and error cases
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
│   │   │   ├── loans.controller.ts
│   │   │   └── dto/
│   │   │       └── create-loan.dto.ts
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
