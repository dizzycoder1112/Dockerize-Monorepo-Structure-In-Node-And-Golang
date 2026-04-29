# Claude Notes

Open-source monorepo template showcasing Dockerized Node + Go services across multiple backend architectures.

## Apps

- `apps/go-layered-server` — 3-layer / layered architecture (handler → service → repository) — refactored from `go-gin-server` template; uses pgx + custom `QueryLogger` tracer, `RepoFactory.UseTransaction` for multi-repo atomicity, `BindJSON` with snake_case validation errors
- `apps/go-ddd-server` — DDD 4-layer (domain / app / infra / interfaces) — extracted from tenderland `go-ops-server`, generic `Order` aggregate
- `apps/ts-restful-api` — Node / TypeScript REST API (Express)
- `apps/ts-grpc-demo` — Node / TypeScript gRPC server consuming `ts-packages/grpc` + `proto/hello.proto` (renamed from `ts-user`; kept as the live consumer of the buf pipeline)

## `go-layered-server` upgrade — done

Refactored from the original `go-gin-server` template up to the structural completeness of tenderland's `go-api-server`, without business logic.

- **C1**: rename — directory, `go.mod` module, import paths, Dockerfile WORKDIR, `package.json` name ✅
- **C2**: structural refactor — all 7 in-scope gaps from CLAUDE.md ✅
  - `pkg/response/` BindJSON + snake_case formatter + InternalServerError caller stash
  - `config/` DSN-style + strict requireEnv (dropped `AA` placeholder)
  - `internal/infra/postgres/` NewPool returns error + QueryLogger tracer (dev=all, prod=slow ≥200ms) + DBConn interface
  - `internal/factory/` composition root with RepoFactory + UseTransaction
  - `internal/router/` health before logger, takes `*Middlewares`
  - `internal/repository/` real CRUD interface + types.go + sentinel errors + TxRunner interface
  - `internal/service/` DealServiceDeps pattern + ctx threading + %w wrapping; `Close` method routed through TxRunner as the extension point for multi-repo atomicity

### Operational notes

- **Default is in-memory** — `go run` works with zero infra. `DATABASE_URL` unset → `factory.NewMemoryRepo()`.
- Set `DATABASE_URL` to flip to the pgx-backed Postgres impl (`factory.NewPostgresRepo(pool)`); failure to ping fatals (intentional — silent fallback would hide config bugs).
- Schema for the postgres path: `deals(id uuid pk, title text, amount bigint, status text, created_at timestamptz, updated_at timestamptz)`. Managed at the monorepo level via Prisma — this app ships no migration.
- `RepoFactory.UseTransaction` degrades to a direct `fn(repos)` call in memory mode — service-level shape stays identical so swapping backends needs no service changes.
- Uses the local `go-packages/logger` via `replace` directive — pattern reusable for any monorepo Go app needing the shared logger.
- `sqlc.yaml` removed (was empty stub with leaked credential — history rewritten via `git filter-repo`, file added to `.gitignore`).

### Out-of-scope (intentional)

- Auth JWT middleware — consumers add their own
- `pkg/utils` Flex* lenient JSON helpers — too business-flavoured
- Integration test harness — it's a template

## Pre-interview smoke test plan (TODO — do tomorrow)

Goal: validate every workspace package the user might `import` mid-interview, so nothing breaks live. Vibe Coding interview is on **2026-04-30**.

### Already validated

- `apps/go-ddd-server` — boots, all 7 routes, in-memory mode ✅
- `apps/go-layered-server` — boots, all 7 routes incl. TxRunner-routed `Close`, in-memory mode ✅
- `apps/ts-grpc-demo` — boots, Connect protocol curl roundtrip ✅
- `pnpm run build` — 8/8 turbo green after `@types/node` centralised at root ✅
- `pnpm run buf:gen` — Go + TS codegen both work, generated Go pb compiles ✅

### Still to validate

User will spin up infra (RabbitMQ, Postgres) locally before the run.

| # | Target | Test |
|---|---|---|
| 1 | `apps/ts-restful-api` | start, hit `/health-check` and one `/api/v1/users` route |
| 2 | `ts-packages/logger` | `node -e` import + emit one line each for info/error/warn/debug |
| 3 | `ts-packages/shared` | import + call one util / read one constant |
| 4 | `ts-packages/grpc` client | `createGreeterClient` against the running ts-grpc-demo — proves client+server pair works |
| 5 | `go-packages/logger` | already exercised by `go-layered-server`; explicit standalone smoke optional |
| 6 | `ts-packages/rabbitMQ` ↔ `go-packages/rabbitMQ` | bring up broker, run **both directions** of producer/consumer — TS→Go and Go→TS. Catches schema/encoding drift across languages |
| 7 | `ts-packages/db` (optional) | needs Postgres; flip `apps/go-layered-server` to `DATABASE_URL=...` to exercise the same connection setup if time permits |
| 8 | gRPC round-trip (stretch) | needs a tiny Go gRPC server stub — currently only `ts-grpc-demo` serves. Skip unless the interview specifically asks for Go-served gRPC |

### Out-of-scope for tomorrow

- README rewrite for `go-layered-server` (currently has the old "Go Counter Server" content) — deferred

## Architecture conventions (target state, both Go apps)

- Composition root lives in `internal/factory/` — one file per layer (`handler.go`, `service.go`, `repository.go`, `middleware.go`)
- `pkg/response/` for unified API response shape + bind-with-validation helper
- `config.Load()` with `getEnv` + sensible fallbacks (templates favour zero-config DX; consumers tighten to `requireEnv` once they own real infra)
- pgx pool with `QueryTracer` (dev logs all queries; prod logs slow only)
- `Middlewares` struct passed into `router.Setup(handlers, middlewares)`
- Sentinel errors defined in domain (DDD) or repository (3-layer) packages; handlers map them to HTTP status
- Per-resource route file under `internal/interfaces/http/router/` (DDD) or `internal/router/` (3-layer)
