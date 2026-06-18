# Library Management System

Distributed library management system built with NestJS and Go microservices.

## Architecture

Two independent services communicating via HTTP:

```
Client
  │
  ▼
library-service (NestJS) ──HTTP──▶ loans-service (Go)
  │                                      │
  ▼                                      ▼
postgres-library                   postgres-loans
```

**library-service** (Port 3000) — main service exposed to the client. Manages books, users and authentication.

**loans-service** (Port 8081) — manages loans. Validates book availability with library-service before registering a loan.

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

```
1. Register user:       POST /users         { name, email, password }
2. Login:              POST /auth/login     { email, password } → access_token
3. Create book (admin): POST /books          { title, author, isbn, year, genre, availableCopies }
4. List books:          GET  /books?author=Borges&page=1&limit=10
5. Create loan:         POST /loans          { userId, bookId }
   → loans-service validates book exists and has copies available via GET /books/:id
   → if valid, registers loan and decrements copies via PATCH /books/:id/copies
6. Return book:         PATCH /loans/:id     { } (no body needed)
   → loans-service updates loan status to "returned"
   → increments available copies via PATCH /books/:id/copies
7. View active loans:   GET  /loans/users/:id
8. View history:        GET  /loans/users/:id/history
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

**Passwords hashed with bcrypt (salt rounds: 10)**
Passwords are never stored or returned in plain text. The `Omit<User, 'password'>` TypeScript type ensures the password field is excluded at the type level from all service responses.

**Pagination with QueryBuilder**
`findAll` uses TypeORM's `createQueryBuilder` to support dynamic filters (author, genre, availability) combined with `skip/take` for pagination. Returns `{ data, total }` so the client can implement pagination UI.

---

## What's pending

- [ ] Role-based guards (admin vs regular user)
- [ ] Tests for library-service (3-4 required)
- [ ] README complete example with curl commands
- [ ] `.env.example` file

## What was intentionally left out

- gRPC between services: HTTP was kept for simplicity and consistency. gRPC would be the natural next step to reduce inter-service latency.
- Rate limiting: out of scope for this assessment.
- Frontend: explicitly excluded by the assessment.
- Kubernetes / service mesh: explicitly excluded by the assessment.

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
│   │   │   └── books.controller.ts
│   │   ├── users/
│   │   │   ├── user.entity.ts
│   │   │   ├── users.module.ts
│   │   │   ├── users.service.ts
│   │   │   └── users.controller.ts
│   │   └── auth/
│   │       ├── auth.module.ts
│   │       ├── auth.service.ts
│   │       ├── auth.controller.ts
│   │       ├── jwt.strategy.ts
│   │       └── jwt-auth.guard.ts
│   └── Dockerfile
└── loans-service/                # Go — loans management
    ├── cmd/api/main.go
    ├── internal/
    │   ├── app.go                # factory / dependency wiring
    │   ├── constants.go          # domain constants
    │   ├── clients/
    │   │   ├── library_client.go # HTTP client for library-service
    │   │   └── errors.go
    │   ├── handler/
    │   │   ├── loan_handler.go
    │   │   ├── routes.go
    │   │   └── errors.go
    │   ├── service/
    │   │   ├── loan_service.go   # BookServiceClient interface defined here
    │   │   ├── library_client.go # concrete implementation
    │   │   └── errors.go
    │   ├── repository/
    │   │   ├── loan_repository.go
    │   │   ├── db.go
    │   │   ├── constants.go      # SQL queries
    │   │   └── errors.go
    │   └── model/
    │       ├── loan.go
    │       └── book.go
    ├── db/schema.sql
    ├── Dockerfile
    └── go.mod
```
