# Flash

[![CI](https://github.com/v9shal/Flash/actions/workflows/ci.yml/badge.svg?branch=feat%2Fauth)](https://github.com/v9shal/Flash/actions/workflows/ci.yml)

Flash is a Go ticketing backend for user authentication, seat discovery, bookings, and payments. PostgreSQL stores booking data; Redis coordinates seat holds and caches seat availability.

## Quick start

Requires Go 1.26.5, PostgreSQL, and Redis. Create a local, untracked `.env` in the repository root with connection details and application secrets:

```dotenv
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/flash_ticket?sslmode=disable
REDISURL=redis://localhost:6379
JWTSECRET=replace-with-a-local-secret
PAYMENT=replace-with-a-local-payment-secret
```

Start the API:

```bash
go run ./cmd/server
```

The service listens on port `8080` by default. PostgreSQL schema setup is not documented as an automated migration in this repository.

## Architecture

```mermaid
flowchart LR
    Client --> API[Gin HTTP API]
    API --> Auth[Authentication]
    API --> Seats[Seat service]
    API --> Bookings[Booking service]
    API --> Payments[Payment service]
    Seats -->|cache and singleflight| Redis[(Redis)]
    Bookings -->|atomic hold script| Redis
    Seats --> DB[(PostgreSQL)]
    Bookings --> DB
    Payments --> DB
    Expiry[Seat expiry worker] --> DB
```

## Design notes

- Seat holds use a Redis script so competing requests cannot both claim the same hold key.
- Seat reads use Redis caching and `singleflight` to coalesce concurrent cache misses for the same event.
- A background worker periodically cancels expired holds; the booking service sets holds to expire after ten minutes.
- The API exposes Prometheus metrics and shuts down its HTTP server on termination signals.

## Limitations and results

- No automated Go tests or benchmark results are currently included in the repository; CI checks build, vet, and test compilation.
- The README's local setup requires an existing PostgreSQL schema; database migrations are not provided as a documented command.
- The checked-in Dockerfile is empty, so the Compose `app` image is not currently buildable. Run the API with Go until container packaging is completed.
