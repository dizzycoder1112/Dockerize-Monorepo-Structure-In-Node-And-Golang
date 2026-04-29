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

## grpc package — pinned to v1 toolchain

connect-es plugin **has no v2** (deprecated after 1.7; the v2 path is "use protoc-gen-es v2 schema only, no separate connect plugin"). Mixing v2 `protoc-gen-es` with v1 `protoc-gen-connect-es` produced incompatible generated files — `_pb` exported `XSchema` but `_connect` imported `X` class → tsc elided the require → `ReferenceError: hello_pb_js_1 is not defined` at runtime.

Resolution: pin everything to v1, matching the working pooktopia template.

- Root `@bufbuild/protoc-gen-es` → `^1.10.1`
- `ts-packages/grpc` runtime: `@bufbuild/protobuf` `^1.10.1`, `@connectrpc/connect`/`connect-node` `^1.7.0`
- src refactored to v1 style: service descriptor imported from `_connect` (not `_pb`), `new HelloReply({...})` instead of `create(HelloReplySchema, ...)`, no `export { create }` from index
- Build switched from `tsc` to `tsup` (centralised at root) for dual ESM+CJS output, matching pooktopia
- ts-grpc-demo `.env` PORT moved to `50051` (gRPC convention) so it doesn't collide with ts-restful-api on `:3000`; `services/index.ts` baseUrl updated to match

### Follow-up: migrate TS side off connect-es to native gRPC (`@grpc/grpc-js` + `ts-proto`)

Current state is "TS server speaks Connect protocol via `connectNodeAdapter`, Go server speaks native gRPC, client picks transport per server". Works, but client has to know which protocol the server speaks — leaky abstraction. Long-term cleaner path is single wire protocol everywhere (native gRPC).

**Why it's worth doing**:
- connect-es 1.x is deprecated (no v2 plugin; v2 direction is "protoc-gen-es v2 schema only without connect plugin"). We're pinned to a dead-end branch.
- `@grpc/grpc-js` is Google's official Node gRPC, actively maintained, same ecosystem as `google.golang.org/grpc`.
- One wire = one set of debugging tools (grpcurl/reflection/wireshark all work everywhere).
- Smaller dependency surface: drops `@bufbuild/protobuf`, `@connectrpc/connect`, `@connectrpc/connect-node`, `protoc-gen-connect-es`, `protoc-gen-es`.

**Scope** (estimated 1–2 hours):
1. `package.json`: drop `@bufbuild/protoc-gen-es` + `@connectrpc/protoc-gen-connect-es`, add `ts-proto`
2. `proto/buf.gen.yaml`: replace `es` + `connect-es` plugins with single `ts_proto`
3. `ts-packages/grpc/package.json`: drop `@bufbuild/protobuf`, `@connectrpc/connect`, `@connectrpc/connect-node`; add `@grpc/grpc-js`
4. `pnpm run buf:gen` — generated TS files structure changes completely (ts-proto outputs idiomatic interfaces instead of class-based messages)
5. Rewrite `ts-packages/grpc/src/{clientFactory,serverFactory,index}.ts` against `@grpc/grpc-js` API
6. Rewrite `apps/ts-grpc-demo/src/handlers/sayHello.handler.ts` (grpc-js promise/callback style)
7. Update `apps/ts-restful-api/src/services/index.ts` + `repositories/user.repository.ts` for new client API
8. Go side unchanged (`protoc-gen-go` + `protoc-gen-go-grpc` already native gRPC)

**Why not now**: too risky pre-interview (2026-04-30). Pick up post-interview.

## go-packages/rabbitMQ — module path renamed

Was `weedza.shop/rabbitmq` (leftover from a different project), renamed to `dizzycoder1112/Dockerize-Monorepo-Structure-In-Node-And-Golang/rabbitmq` to match the rest of the monorepo's module naming convention. No app consumes this package yet, so the only file affected was `go-packages/rabbitMQ/go.mod`.

## Pre-interview smoke test plan (TODO — do tomorrow)

Goal: validate every workspace package the user might `import` mid-interview, so nothing breaks live. Vibe Coding interview is on **2026-04-30**.

### Already validated

- `apps/go-ddd-server` — boots, all 7 routes, in-memory mode ✅
- `apps/go-layered-server` — boots, all 7 routes incl. TxRunner-routed `Close`, in-memory mode ✅
- `apps/ts-grpc-demo` — boots, Connect protocol curl roundtrip ✅
- `pnpm run build` — 8/8 turbo green after `@types/node` centralised at root ✅
- `pnpm run buf:gen` — Go + TS codegen both work, generated Go pb compiles ✅

### Still to validate

User will spin up infra (RabbitMQ, Postgres) locally before the run. Order: zero-infra first, infra-dependent last.

| # | Status | Target | Test |
|---|---|---|---|
| 1 | ✅ | `ts-packages/logger` | `node -e` import + emit one line each for info/warn/error/debug — verified pretty transport, ISO timestamp, serviceName, context all working |
| 2 | ✅ | `ts-packages/shared` | `node -e` import `./constants` + `./utils` — verified `SERVICE_NAME` enum (3 values) + `sleep(150)` resolved in 152ms |
| 3 | ✅ | `apps/ts-restful-api` | boots on :3000, `/health-check` → 200 OK; `/api/v1/users/sayHello?name=X` → e2e through grpc to ts-grpc-demo → `{"message":"You said X"}` |
| 4 | ✅ | `ts-packages/grpc` client ↔ `ts-grpc-demo` | done together with #3 — full client→server roundtrip via Connect-over-H2 on :50051 |
| 5 | ✅ | `ts-packages/rabbitMQ` ↔ `go-packages/rabbitMQ` | broker on OrbStack `rabbitmq.rabbitmq.orb.local:5672`, both directions verified — TS→Go: `GO_CONSUMED: {"hello":"from ts",...}` / Go→TS: `TS_CONSUMED: {"hello":"from go",...}`. JSON round-trip clean across both `amqplib`(TS) and `amqp091-go`(Go). Smoke runner lives at `tmp/rabbit-smoke/` (gitignored) |
| 6 | ⏳ | `ts-packages/db` (optional) | needs Postgres; flip `apps/go-layered-server` to `DATABASE_URL=...` to exercise the same connection setup if time permits |

### Skipped (intentional)

- `go-packages/logger` standalone — already exercised by `go-layered-server`, explicit smoke not needed

### Stretch — done

- ✅ `apps/go-grpc-demo` on :50052 — minimal Go server using `go-packages/grpc.NewServer()` wrapper, registers `Greeter`, EnableReflection. Verified single TS `createGreeterClient` factory hits **both** ts-grpc-demo (:50051) and go-grpc-demo (:50052) with identical API — `connectNodeAdapter` routes `application/grpc` to native handler so client speaks one protocol everywhere. `clientFactory.ts` now always uses `createGrpcTransport` (no protocol option — was a leaky abstraction)

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
