# Library Management System

Distributed library management system — two independent microservices communicating over HTTP.

| Service | Stack | Port | Database |
|---------|-------|------|----------|
| **library-service** | NestJS + TypeORM | 3000 | postgres-library (5432) |
| **loans-service** | Go + Chi + pgx | 8081 (internal) | postgres-loans (5433) |

---

## Architecture

```
Client
  │
  ▼  (all client traffic)
library-service  (NestJS, :3000)
  │  ├─ books, users, auth → postgres-library
  │  └─ proxies loan ops ──────HTTP──────▶ loans-service (Go, :8081)
  │                                              └─ loans → postgres-loans
```

**library-service** is the single entry point for all clients. It manages books, users, and authentication, and proxies loan operations to loans-service.

**loans-service** is internal — not exposed directly to clients. It manages loan records and validates book availability by calling library-service on each create/return.

Each service owns its own PostgreSQL database, enforcing true independence at the data layer.

**Inter-service auth:** loans-service generates a short-lived JWT (`role: "service"`, 1-min expiry) signed with the shared `JWT_SECRET` on every outbound call. library-service validates it with the standard `JwtAuthGuard` — no new auth mechanism needed.

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

Starts 4 containers: library-service, loans-service, postgres-library, postgres-loans.

loans-service waits for postgres-loans to be healthy before starting (Docker healthcheck). The loans schema is applied automatically on first run from `loans-service/db/schema.sql`.

> **Data persistence:** Use `docker compose down` (without `-v`) to stop while keeping data. `docker compose down -v` resets everything.

### Environment variables

Copy `.env.example` to `.env` before running locally without Docker.

```env
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/library_db
LOANS_SERVICE_URL=http://localhost:8081
LIBRARY_SERVICE_URL=http://localhost:3000
JWT_SECRET=your_secret_here
PORT=3000
```

> **Note:** Both services must share the same `JWT_SECRET`. Never hardcode it in production — use a secrets manager.

---

## API Reference

All client requests go through **library-service** on port 3000. Swagger UI is available at `http://localhost:3000/api` when the stack is running.

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
| PATCH | `/books/:id/copies` | JWT (internal) | Update available copies — called by loans-service only |

**Filters for `GET /books`:** `?author=&genre=&available=true&page=1&limit=10`

### Loans

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/loans` | JWT | Create loan (proxied to loans-service) |
| PATCH | `/loans/:id` | JWT | Return a book (proxied to loans-service) |
| GET | `/loans/users/:userId` | JWT | Active loans for a user |
| GET | `/loans/users/:userId/history` | JWT | Full loan history for a user |

### Health

| Method | Endpoint | Service |
|--------|----------|---------|
| GET `/health` | port 3000 | library-service |
| GET `/health` | port 8081 | loans-service (direct) |

---

## Role Permissions

| Action | admin | user |
|--------|-------|------|
| Register / Login | ✅ | ✅ |
| List / Get books | ✅ | ✅ |
| Create / Update / Delete book | ✅ | ❌ |
| List all users | ✅ | ❌ |
| Get / Update own user | ✅ | ✅ |
| Change a user's role | ✅ | ❌ |
| Delete user | ✅ | ❌ |
| Create / Return loan | ✅ | ✅ |
| View active loans / history | ✅ | ✅ |

> Loan endpoints require authentication but do not restrict by role — any authenticated user can borrow and return books. Suggested improvement: restrict `GET /loans/users/:userId` so regular users can only query their own loans (validate `:userId` matches the JWT `sub`).

---

## Complete Flow Example

```bash
# 1. Login as admin
TOKEN=$(curl -s -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@library.com","password":"adminpass"}' \
  | jq -r '.access_token')

# 2. Check a book (no auth required)
curl -s http://localhost:3000/books/1 | jq '{id,title,author,available_copies}'
# {"id":1,"title":"The Go Programming Language","author":"Alan Donovan","available_copies":3}

# 3. Create a loan
curl -s -X POST http://localhost:3000/loans \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"book_id":1}' | jq .
# {"id":1,"user_id":1,"book_id":1,"loaned_at":"...","status":"active"}

# 4. Verify copies decremented
curl -s http://localhost:3000/books/1 | jq .available_copies
# 2  (was 3)

# 5. Return the book
curl -s -X PATCH http://localhost:3000/loans/1 \
  -H "Authorization: Bearer $TOKEN" | jq .
# {"id":1,...,"returned_at":"...","status":"returned"}

# 6. Verify copies restored
curl -s http://localhost:3000/books/1 | jq .available_copies
# 3  (restored)

# 7. View loan history
curl -s http://localhost:3000/loans/users/1/history \
  -H "Authorization: Bearer $TOKEN" | jq .
```

---

## End-to-End Testing

All tests run against `http://localhost:3000` after `docker compose up --build`.

### 1. Seed the database

```bash
chmod +x scripts/seed.sh && ./scripts/seed.sh
```

Creates 2 admin users, 1 regular user, 4 books, and 3 loans.

### 2. Verify database state

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

### 3. CRUD and validation tests

```bash
chmod +x scripts/test_update_delete.sh && ./scripts/test_update_delete.sh
```

Tests UPDATE and DELETE for books and users, covering validation edge cases (invalid fields, missing auth, wrong role, not found). Expected responses captured in `scripts/test_results_update_delete.txt`.

### 4. Full functional tests

```bash
chmod +x scripts/test_full.sh && ./scripts/test_full.sh
```

Tests the complete loan lifecycle and error cases. Expected highlights:
- Loan creation → `201` with `status: "active"`
- Available copies decrement after loan, increment after return
- Second loan for same user + book → `409 Conflict`
- Exhausting all copies → `409 Conflict`
- Missing token → `401 Unauthorized`
- Non-admin creating a book → `403 Forbidden`

Captured expected responses in `scripts/test_results_full.txt`.

---

## Development

### library-service (NestJS)

```bash
cd library-service
npm ci                  # install deps
npm run start:dev       # dev mode with watch
npm test                # unit tests
npm run test:cov        # with coverage
npm run test:e2e        # e2e tests
npm run lint            # lint + fix
npm run build           # compile to dist/

# Run a single test file
npx jest src/books/books.service.spec.ts
```

### loans-service (Go)

```bash
cd loans-service
go test ./internal/...               # all unit tests
go test ./internal/service/...       # service layer only
go test ./internal/handler/...       # handler layer only
go build ./cmd/api/...               # compile
go vet ./...                         # lint
```

---

## CI

GitHub Actions runs unit tests for both services on every push and pull request.

```
.github/workflows/ci.yml
├── test-library-service   → npm ci && npm test  (Node.js 20)
└── test-loans-service     → go test ./internal/...  (Go 1.21)
```

No database or Docker required — all tests use mocks.

---

## Technical Decisions

### loans-service (Go)

**HTTP Framework: Chi**
Chi is closer to the standard library and more idiomatic Go than Gin, which adds routing magic that hides what's happening. Since testability is a goal, Chi was the better fit.

**Database driver: pgx over GORM**
`pgx` with `database/sql` shows explicit control over queries. GORM can hide performance problems behind abstraction — in a technical assessment, raw SQL communicates more than an ORM.

**Two separate PostgreSQL instances**
Each service owns its data completely. A shared database creates coupling at the data layer — if library-service changes its schema, loans-service breaks. Two instances enforce true independence. Trade-off: higher resource usage. Justified for a microservices architecture.

**Multi-stage Docker build**
Final image contains only the compiled binary + Alpine Linux (~5MB + binary). No Go compiler, no source code.

**Project structure and constants**
- SQL queries live in `repository/constants.go`
- Domain status values (`active`, `returned`) and copy actions (`increment`, `decrement`) live in `service/constants.go`
- Errors are declared per-layer in `errors.go` files (`repository`, `service`, `handler`, `clients`)

**Dependency inversion on BookServiceClient**
The service layer defines the interface it needs (`BookService`); the concrete implementation (`clients/library_client.go`) satisfies it via Go duck typing. This keeps the service layer independent of HTTP details and fully testable with mocks.

**Factory pattern (`internal/app.go`)**
All infrastructure (pgx pool, repository, service) is constructed in one place. `cmd/api/main.go` only wires the HTTP layer. The factory deliberately excludes the HTTP framework to avoid coupling infrastructure with transport.

**Schema: partial unique index**
```sql
CREATE UNIQUE INDEX ON loans(user_id, book_id) WHERE status = 'active';
```
Prevents duplicate active loans for the same user+book while allowing re-borrow after return. A simple UNIQUE constraint would block that.

**`PATCH /loans/:id` for returns**
Instead of a verb-in-URL pattern (`/return`), `PATCH /loans/:id` represents a state transition. The service layer method `BookReturned` always sets status to `returned` — the status is a consequence of the operation, not an input from the client.

**Inter-service auth: short-lived JWT**
loans-service generates a JWT (`role: "service"`, 1-min expiry) signed with the shared `JWT_SECRET` on each outbound call. library-service validates it with the standard `JwtAuthGuard`. Short expiry limits blast radius of a token leaked in transit.

**503 on library-service unavailable**
Network errors from `http.Client.Do()` are wrapped as `ErrLibraryServiceUnavailable`, unwrapped at the handler layer with `errors.Is()`, and returned as clean `503 Service Unavailable`.

**Distributed transaction handling**
Loan operations span two systems (loans DB + library-service copy count). Best-effort consistency via compensating calls:

| Operation | Step 1 | Step 2 | On step 2 failure |
|-----------|--------|--------|-------------------|
| Create loan | Decrement copies | Insert loan record | Compensating increment |
| Return book | Increment copies | Update loan status | Compensating decrement |

A failure in the compensating call leaves systems out of sync — a known trade-off. Production would use an outbox pattern or saga coordinator.

---

### library-service (NestJS)

**TypeORM over Prisma**
TypeORM has more mature complex relational mappings and better TypeScript integration with NestJS decorators. Both were valid — TypeORM was selected for stability.

**`autoLoadEntities: true`**
Each module registers its own entity with `forFeature()` instead of listing all entities in `TypeOrmModule.forRoot()`. Scales better as the entity count grows.

**`synchronize: true`**
TypeORM auto-creates/updates tables from entity definitions. Appropriate for development and this assessment. In production: use `typeorm migration:generate` + `typeorm migration:run`.

**JWT with Passport**
`@nestjs/passport` + `passport-jwt` is the standard NestJS auth pattern. `JwtStrategy` validates the `Authorization: Bearer` header. `JwtAuthGuard` is applied per-endpoint with `@UseGuards()`.

**Role-based access: custom decorator + guard**
`@Roles('admin')` stores metadata on the endpoint. `RolesGuard` reads it via `Reflector` and compares against the JWT payload role. Endpoints without `@Roles` are accessible to any authenticated user.

**Input validation: class-validator**
All endpoints use typed DTOs with `class-validator` decorators. Global `ValidationPipe` (`whitelist + forbidNonWhitelisted`) in `main.ts` rejects unknown fields at the boundary — including the loans proxy, which validates `user_id` and `book_id` before forwarding.

**Atomic copy update**
`updateCopies` uses a single `UPDATE ... SET available_copies = available_copies + $delta WHERE available_copies + $delta >= 0 RETURNING *`. Eliminates the race condition where two concurrent loan requests could both read the same copy count.

**Passwords hashed with bcrypt (salt rounds: 10)**
Never stored or returned in plain text. `Omit<User, 'password'>` excludes the field at the type level from all service responses.

**Pagination with QueryBuilder**
`findAll` uses TypeORM's `createQueryBuilder` for dynamic filters (author, genre, availability) combined with `skip/take`. Returns `{ data, total }` for client-side pagination.

---

## Testing

| Service | Command | What's covered |
|---------|---------|----------------|
| loans-service | `go test ./internal/service/...` | CreateLoan happy path, book not available, BookReturned, GetActiveLoans |
| loans-service | `go test ./internal/handler/...` | CreateLoan 201, invalid body 400 |
| library-service | `npm test` | BooksService (create, findOne, updateCopies), AuthService (valid login, user not found, wrong password) |

All tests use mocks — no running database required. Tests also run automatically on every push via GitHub Actions.

---

## Project Structure

```
library-system/
├── .github/workflows/ci.yml          # GitHub Actions — tests on push/PR
├── docker-compose.yml
├── .env.example
├── scripts/
│   ├── seed.sh                        # creates test users, books, loans
│   ├── test_update_delete.sh          # CRUD + validation tests
│   └── test_full.sh                   # auth, loans lifecycle, error cases
│
├── library-service/                   # NestJS — books, users, auth, loans proxy
│   └── src/
│       ├── auth/                      # JWT strategy, guards, decorators
│       ├── books/                     # CRUD + atomic copy update
│       ├── users/                     # CRUD + bcrypt
│       └── loans/                     # thin HTTP proxy to loans-service
│
└── loans-service/                     # Go — loans management
    ├── cmd/api/main.go                # HTTP wiring only
    ├── internal/
    │   ├── app.go                     # factory: pool → repository → service
    │   ├── clients/                   # HTTP client for library-service
    │   ├── handler/                   # Chi handlers, error mapping
    │   ├── service/                   # business logic, interfaces
    │   ├── repository/                # pgx queries, SQL constants
    │   └── model/                     # Loan, Book, request types
    └── db/schema.sql                  # applied on first container boot
```

---

## What Was Left Out / Known Limitations

### Conscious design decisions

| Item | Decision & Reasoning | Risk / Production Fix |
|------|------------------------|--------------------------|
| Cross-service consistency on loan creation/return | No real ACID transaction is possible across two separate databases (loans-service and library-service), so the flow was inverted: copies are updated in library-service first, then the loan is created/updated in loans-service. If the second step fails, a compensating call reverts the copies update — a manually orchestrated Saga with compensating action. The compensating call uses `context.Background()` instead of the original request context, since the original context could already be cancelled/expired and the rollback still needs to run | If the compensating call itself fails, the system is left inconsistent with no automatic retry — currently mitigated only by logging for manual reconciliation. Production fix: outbox pattern (persist the compensation as an event in loans-service's local transaction, retried by a background worker until success) |
| TypeORM `synchronize: true` | Used instead of migrations for faster iteration within the time-boxed evaluation | Production: `migration:generate` + `migration:run` with versioned migration files |
| HTTP instead of gRPC | Kept for simplicity | gRPC would be the natural next step for lower latency and strongly-typed contracts via protobuf |
| Unit tests over integration tests | Integration tests require a running database; business logic was covered with unit tests instead, given the time constraint | One exception, see race condition below |
| Rate limiting / Frontend / Kubernetes | Out of scope — rate limiting was an optional bonus item; frontend/Kubernetes explicitly excluded by the assignment | — |

### Known limitations (found during testing, not resolved due to time)

| Item | What happens | Why it wasn't resolved | Production fix |
|------|----------------|--------------------------|--------------------|
| TOCTOU race condition on available copies | Two concurrent loan requests for the same book with 1 copy left can both pass the `copiasDisponibles > 0` check before either decrements the counter | Found during manual testing, not in the original plan. Not verified with `go test -race` — Go's race detector catches memory races between goroutines in one process, not this kind of database-level race across separate HTTP requests. Correct verification requires a concurrency test (N parallel requests, assert only 1 succeeds), not implemented before submission due to time; planned as a follow-up | Atomic conditional update: `UPDATE books SET copias_disponibles = copias_disponibles - 1 WHERE id = $1 AND copias_disponibles > 0`, checking `RowsAffected() == 0` to reject when no copies are available |
| Loan ownership restriction | `GET /loans/users/:userId` doesn't validate that the caller is that user or an admin — any authenticated user can list another user's loans | Not an explicit requirement for Servicio B (it has no auth of its own — that's Servicio A's responsibility), but surfaced as a gap in manual testing; no time to add the authorization check | Servicio A should derive `userId` from the JWT instead of trusting the path parameter, or validate caller identity/role before proxying to loans-service |
