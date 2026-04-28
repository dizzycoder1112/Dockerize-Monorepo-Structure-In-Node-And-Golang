# Claude Notes

Open-source monorepo template showcasing Dockerized Node + Go services across multiple backend architectures.

## Apps

- `apps/go-hello` — minimal Go hello-world
- `apps/go-layered-server` — 3-layer / layered architecture (handler → service → repository) — **C1 rename done; C2 refactor pending, see below**
- `apps/go-ddd-server` — DDD 4-layer (domain / app / infra / interfaces) — extracted from tenderland `go-ops-server`, generic `Order` aggregate
- `apps/ts-restful-api`, `apps/ts-user` — Node / TypeScript apps

## In-flight: `go-layered-server` upgrade

Goal: bring the 3-layer template up to the completeness of tenderland's `go-api-server`, without copying business logic.

### Status

- [x] **C1 rename complete** — `go-gin-server` → `go-layered-server` (folder by user, module/imports/Dockerfile WORKDIR/package.json by Claude). Build green.
- [ ] **C2 refactor** — gaps #1–7 below, scope confirmed

### Commit plan (B route — rename first, refactor second)

- **C1**: pure rename — directory, `go.mod` module, all import paths, `Dockerfile` / `Dockerfile.local` WORKDIR, `package.json` name ✅
- **C2**: refactor (gaps #1–7 below)

### In-scope gaps

| # | Gap | Effort |
|---|---|---|
| 1 | `pkg/response/` — `BindJSON`, unified response shape, sentinel errors (`ErrNotFound`) | S |
| 2 | `config/` — switch to DSN-style `DATABASE_URL`, `requireEnv`, env-gated branches, drop `AA` placeholder | S |
| 3 | `db/` → `infra/postgres/` — `NewPool` returns error (no `MustNew`), pgx `QueryTracer` (dev=all, prod=slow) | M |
| 4 | `factory/` composition root — extract wiring out of `main.go` | M |
| 5 | `router/` — split per resource, `Middlewares` struct, health route before logger | S |
| 6 | `repository/` — flesh out empty interface, add `types.go`, `TxRunner` for multi-repo transactions | M |
| 7 | `service/` — Deps-struct pattern, `ctx` threading, `%w` error wrapping | S |

### Out-of-scope (intentional)

- Auth JWT middleware — consumers add their own auth
- `pkg/utils` `FlexInt32` / `FlexFloat64` lenient JSON helpers — too business-flavoured
- Integration test harness — this is a template

### Bonus fix during C1

- `Dockerfile` and `Dockerfile.local` both have `WORKDIR /app/apps/go-counter-server` (broken — directory is `go-gin-server`). Fix while renaming.

## Architecture conventions (target state, both Go apps)

- Composition root lives in `internal/factory/` — one file per layer (`handler.go`, `service.go`, `repository.go`, `middleware.go`)
- `pkg/response/` for unified API response shape + bind-with-validation helper
- `config.Load()` with `requireEnv` + env-gated logging
- pgx pool with `QueryTracer` (dev logs all queries; prod logs slow only)
- `Middlewares` struct passed into `router.Setup(handlers, middlewares)`
- Sentinel errors defined in domain (DDD) or repository (3-layer) packages; handlers map them to HTTP status
- Per-resource route file under `internal/interfaces/http/router/` (DDD) or `internal/router/` (3-layer)
