# library-service

NestJS service — manages books, users, authentication, and proxies loan operations to loans-service.

## Run locally

```bash
npm install
npm run start:dev
```

## Run tests

```bash
npm test
```

## Environment variables

| Variable | Description |
|----------|-------------|
| `DATABASE_URL` | PostgreSQL connection string |
| `LOANS_SERVICE_URL` | Base URL of loans-service (e.g. `http://localhost:8081`) |
| `JWT_SECRET` | Shared secret for JWT signing — must match loans-service |
| `PORT` | Port to listen on (default: `3000`) |

See `.env.example` in the root of the repo.
