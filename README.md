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

The loans schema is applied automatically on first run from `loans-service/db/schema.sql`.

### Environment variables

Copy `.env.example` to `.env` and fill in the values before running locally without Docker.

```env
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/library_db
LOANS_SERVICE_URL=http://localhost:8081
JWT_SECRET=your_secret_here
LIBRARY_SERVICE_URL=http://localhost:3000
```

> **Note:** JWT_SECRET must never be hardcoded in production. Use a secrets manager or external environment variable injection.

---

## Example Flow

### 1. Register a user

```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Verónica", "email": "vero@test.com", "password": "secret123", "role": "admin"}'
```

### 2. Login

```bash
curl -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "vero@test.com", "password": "secret123"}'
```

Response:
```json
{ "access_token": "eyJhbGci..." }
```

### 3. Create a book (admin only)

```bash
curl -X POST http://localhost:3000/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"title": "Clean Code", "author": "Robert Martin", "isbn": "9780132350884", "year": 2008, "genre": "tech", "availableCopies": 3}'
```

### 4. List books with filters and pagination

```bash
curl "http://localhost:3000/books?author=Martin&genre=tech&available=true&page=1&limit=10"
```

### 5. Create a loan

```bash
curl -X POST http://localhost:3000/loans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"user_id": 1, "book_id": 1}'
```

library-service forwards the request to loans-service, which validates book availability with library-service before persisting the loan.

### 6. Return a book

```bash
curl -X PATCH http://localhost:3000/loans/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

loans-service updates the loan status to `returned` and increments the available copies in library-service.

### 7. View active loans for a user

```bash
curl http://localhost:3000/loans/users/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 8. View loan history

```bash
curl http://localhost:3000/loans/users/1/history \
  -H "Authorization: Bearer YOUR_TOKEN"
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

---

### library-service (NestJS)

**TypeORM over Prisma**
TypeORM is more mature for complex relational mappings and has better TypeScript integration with NestJS decorators. Both were valid choices — TypeORM was selected for stability.

**autoLoadEntities: true**
Instead of listing every entity in `TypeOrmModule.forRoot()`, each module registers its own entity with `forFeature()`. Scales better as the number of entities grows.

**synchronize: true**
TypeORM automatically creates/updates tables based on entity definitions. Appropriate for development and this assessment. In production, use migrations.

**JWT with Passport**
`@nestjs/passport` + `passport-jwt` is the standard NestJS authentication pattern. The `JwtStrategy` validates the token from the `Authorization: Bearer` header. `JwtAuthGuard` is applied per-endpoint with `@UseGuards()`.

**Role-based access control with custom decorator + guard**
`@Roles('admin')` decorator stores metadata on the endpoint. `RolesGuard` reads that metadata via `Reflector` and compares it against the role in the JWT payload. Endpoints without `@Roles` are accessible to any authenticated user.

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
- **User existence validation in loans-service:** loans-service does not call library-service to verify that the userId is a valid user before creating a loan. A dedicated `/users/:id` validation call could be added, but adds latency for the common path.
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
├── library-service/              # NestJS — books, users, auth
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