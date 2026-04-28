# Claude Notes

Open-source monorepo template showcasing Dockerized Node + Go services across multiple backend architectures.

## Apps

- `apps/go-hello` — minimal Go hello-world
- `apps/go-layered-server` — 3-layer / layered architecture (handler → service → repository) — refactored from `go-gin-server` template; uses pgx + custom `QueryLogger` tracer, `RepoFactory.UseTransaction` for multi-repo atomicity, `BindJSON` with snake_case validation errors
- `apps/go-ddd-server` — DDD 4-layer (domain / app / infra / interfaces) — extracted from tenderland `go-ops-server`, generic `Order` aggregate
- `apps/ts-restful-api`, `apps/ts-user` — Node / TypeScript apps

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

- Requires Postgres + a `deals(id uuid pk, title text, amount bigint, status text, created_at timestamptz, updated_at timestamptz)` table. Schema is managed via Prisma at the monorepo level — this app does not ship its own migration. `go run` will fatal on Ping if `DATABASE_URL` is unreachable.
- Uses the local `go-packages/logger` via `replace` directive — pattern reusable for any monorepo Go app needing the shared logger.
- `sqlc.yaml` removed (was empty stub with leaked credential — history rewritten via `git filter-repo`, file added to `.gitignore`).

### Out-of-scope (intentional)

- Auth JWT middleware — consumers add their own
- `pkg/utils` Flex* lenient JSON helpers — too business-flavoured
- Integration test harness — it's a template
- README rewrite — original "Go Counter Server" content still present, deferred

## Architecture conventions (target state, both Go apps)

- Composition root lives in `internal/factory/` — one file per layer (`handler.go`, `service.go`, `repository.go`, `middleware.go`)
- `pkg/response/` for unified API response shape + bind-with-validation helper
- `config.Load()` with `getEnv` + sensible fallbacks (templates favour zero-config DX; consumers tighten to `requireEnv` once they own real infra)
- pgx pool with `QueryTracer` (dev logs all queries; prod logs slow only)
- `Middlewares` struct passed into `router.Setup(handlers, middlewares)`
- Sentinel errors defined in domain (DDD) or repository (3-layer) packages; handlers map them to HTTP status
- Per-resource route file under `internal/interfaces/http/router/` (DDD) or `internal/router/` (3-layer)
