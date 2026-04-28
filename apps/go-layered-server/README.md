# go-layered-server

A 3-layer (handler → service → repository) Go HTTP server built on Gin. Pairs with `go-ddd-server` to showcase two backend architectures inside the same monorepo.

## Quickstart

```bash
go run ./cmd/main.go
# server on :8080, in-memory repositories
```

Set `DATABASE_URL` to flip to the pgx-backed Postgres backend:

```bash
DATABASE_URL=postgres://postgres:postgres@localhost:5432/dev?sslmode=disable go run ./cmd/main.go
```

Schema (when using Postgres):

```sql
CREATE TABLE deals (
  id          UUID PRIMARY KEY,
  title       TEXT        NOT NULL,
  amount      BIGINT      NOT NULL,
  status      TEXT        NOT NULL DEFAULT 'open',
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

## Layout

```
apps/go-layered-server/
├── cmd/main.go                    composition root
├── internal/
│   ├── config/                    env loader (getEnv with fallbacks)
│   ├── infra/postgres/            NewPool + QueryLogger tracer + DBConn
│   ├── repository/
│   │   ├── interfaces.go          DealRepository, Repositories, TxRunner
│   │   ├── types.go               DealRow normalized struct
│   │   ├── errors.go              sentinel errors (ErrNotFound, ...)
│   │   ├── memory/                in-memory backend (default)
│   │   └── postgres/              pgx-backed backend
│   ├── service/                   DealServiceDeps pattern, ctx threading
│   ├── handler/                   gin handlers, BindJSON + sentinel mapping
│   ├── middleware/                cors, logger, Middlewares struct
│   ├── router/                    per-resource route files
│   └── factory/                   composition root: NewMemoryRepo / NewPostgresRepo / NewService / NewHandler / NewMiddleware
└── pkg/response/                  unified API response shape + BindJSON
```

## API

| Method | Path                          | Notes                            |
|--------|-------------------------------|----------------------------------|
| GET    | `/health`                     | liveness probe                   |
| GET    | `/api/v1/deals`               | list                             |
| GET    | `/api/v1/deals/:id`           |                                  |
| POST   | `/api/v1/deals`               | `{title, amount}`                |
| PATCH  | `/api/v1/deals/:id`           | `{title, amount}`                |
| POST   | `/api/v1/deals/:id/close`     | runs through `TxRunner`          |
| DELETE | `/api/v1/deals/:id`           |                                  |

Validation errors return `400` with snake_cased field names. Unknown IDs return `404`. Invalid status transitions / args return `400`/`409`.

## Notes

- `RepoFactory.UseTransaction` degrades to a direct `fn(repos)` call in memory mode — service code is identical across both backends.
- `QueryLogger` (the pgx tracer in `infra/postgres/tracer.go`) logs every query in dev/staging at Debug, only slow queries (≥200 ms) at Warn in production.
- Composition root lives in `internal/factory/` — one file per layer. Add new services / repos by extending these factories rather than editing `cmd/main.go`.
- Logger comes from `go-packages/logger` via a `replace` directive — same pattern any monorepo Go app can copy.
