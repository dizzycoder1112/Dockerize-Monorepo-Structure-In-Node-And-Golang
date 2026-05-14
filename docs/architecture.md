# Architecture

## Apps

| App | Path | Architecture | Notes |
|---|---|---|---|
| `go-layered-server` | `apps/go-layered-server` | 3-layer (handler → service → repository) | pgx + custom `QueryLogger` tracer, `RepoFactory.UseTransaction` for multi-repo atomicity, `BindJSON` with snake_case validation errors. Refactored from `go-gin-server` template. |
| `go-ddd-server` | `apps/go-ddd-server` | DDD 4-layer (domain / app / infra / interfaces) | Generic `Order` aggregate. Extracted from tenderland `go-ops-server`. |
| `go-grpc-demo` | `apps/go-grpc-demo` (`:50052`) | Minimal native-gRPC Go server | Uses `go-packages/grpc.NewServer()` wrapper, registers `Greeter`, `EnableReflection`. |
| `ts-restful-api` | `apps/ts-restful-api` (`:3000`) | Express REST API | `/health-check`, `/api/v1/users/sayHello` (proxies to `ts-grpc-demo` via the shared gRPC client). |
| `ts-grpc-demo` | `apps/ts-grpc-demo` (`:50051`) | gRPC server | Consumes `ts-packages/grpc` + `proto/hello.proto`. Speaks Connect protocol via `connectNodeAdapter`, which routes `application/grpc` traffic to a native handler. |

## Shared packages

| Package | Purpose |
|---|---|
| `ts-packages/logger` | Pino-based logger with pretty transport, ISO timestamp, serviceName, context. |
| `ts-packages/shared` | `SERVICE_NAME` enum + utility helpers (`sleep`, etc.). |
| `ts-packages/grpc` | Connect-RPC client + server factories built on `proto/hello.proto` codegen. Client uses `createGrpcTransport` for native gRPC wire. |
| `ts-packages/rabbitMQ` | TS RabbitMQ wrapper (`amqplib`). Counterpart of `go-packages/rabbitMQ`. |
| `ts-packages/db` | Prisma schema at `prisma/main/schema.prisma`; runtime via Kysely with `CamelCasePlugin` for auto camelCase↔snake_case. |
| `go-packages/logger` | Shared Go logger; consumed by `go-layered-server` via a `replace` directive. |
| `go-packages/grpc` | `NewServer()` wrapper for native gRPC + reflection. |
| `go-packages/rabbitMQ` | Go RabbitMQ wrapper (`amqp091-go`). No app consumes it yet. |

## Conventions

### Both Go apps (target state)

- Composition root lives in `internal/factory/` — one file per layer (`handler.go`, `service.go`, `repository.go`, `middleware.go`).
- `pkg/response/` for unified API response shape + bind-with-validation helper.
- `config.Load()` with `getEnv` + sensible fallbacks (templates favour zero-config DX; consumers tighten to `requireEnv` once they own real infra).
- pgx pool with `QueryTracer` — dev logs all queries; prod logs slow only (≥ 200 ms).
- `Middlewares` struct passed into `router.Setup(handlers, middlewares)`.
- Sentinel errors defined in domain (DDD) or repository (3-layer) packages; handlers map them to HTTP status.
- Per-resource route file under `internal/interfaces/http/router/` (DDD) or `internal/router/` (3-layer).

### `go-packages/*` module path

Go packages under `go-packages/*` use `dizzycoder1112/Dockerize-Monorepo-Structure-In-Node-And-Golang/<pkg>` as their module path, to match the rest of the monorepo's naming.

## `go-layered-server` runtime

- **Default is in-memory** — `go run` works with zero infra. `DATABASE_URL` unset → `factory.NewMemoryRepo()`.
- Set `DATABASE_URL` to flip to the pgx-backed Postgres impl (`factory.NewPostgresRepo(pool)`). Failure to ping fatals (intentional — silent fallback would hide config bugs).
- Schema (Postgres path): `deals(id uuid pk, title text, amount bigint, status text, created_at timestamptz, updated_at timestamptz)`. Managed at the monorepo level via Prisma — this app ships no migration.
- `RepoFactory.UseTransaction` degrades to a direct `fn(repos)` call in memory mode, so service-level shape stays identical across backends.
- The `Close` service method routes through `RepoFactory.UseTransaction` even though it touches a single repo. `UseTransaction` is the extension point for multi-repo atomicity — routing single-repo writes through it now means future deals-touching-other-resources logic ports cleanly without changing service signatures.

## Intentional exclusions (`go-layered-server`)

Deliberately left to consumers (it's a template):

- Auth JWT middleware — consumers add their own.
- `pkg/utils` Flex* lenient JSON helpers — too business-flavoured for a template.
- Integration test harness.

## Key decisions

### TS grpc toolchain pinned to v1

**Decision**: TS-side grpc tooling pinned to v1 of `@bufbuild/protoc-gen-es`, `@bufbuild/protobuf`, `@connectrpc/connect`, `@connectrpc/protoc-gen-connect-es`.

**Why**: `protoc-gen-connect-es` has no v2 release — connect-es 1.7 is the last version, and the v2 direction is "use `protoc-gen-es` v2 schema-only, no separate connect plugin". Mixing v2 `protoc-gen-es` with v1 `protoc-gen-connect-es` produces incompatible generated files: `_pb` exports `XSchema` while `_connect` imports the `X` class → tsc elides the `require` → `ReferenceError: hello_pb_js_1 is not defined` at runtime.

**Current state**:
- Root `@bufbuild/protoc-gen-es` `^1.10.1`.
- `ts-packages/grpc` runtime: `@bufbuild/protobuf` `^1.10.1`, `@connectrpc/connect` / `connect-node` `^1.7.0`.
- v1 codegen style: service descriptor imported from `_connect` (not `_pb`); `new HelloReply({...})` instead of `create(HelloReplySchema, ...)`.
- Build via `tsup` (centralised at root) for dual ESM+CJS output.

**Future**: tracked as the draft plan `docs/plans/active/ts-grpc-native-migration.md` — migrate off connect-es onto native `@grpc/grpc-js` + `ts-proto`.

### Unified gRPC client transport

**Decision**: `ts-packages/grpc/clientFactory.ts` always uses `createGrpcTransport`, with no per-server `protocol` option.

**Why**: on the TS server side, `connectNodeAdapter` routes `application/grpc` traffic to native handlers, so a single client factory can hit both the TS server (`:50051`) and the Go server (`:50052`) with the same API. A per-server `protocol` switch was a leaky abstraction — the client should not need to know which protocol the server speaks.

### Build dual ESM+CJS via root `tsup`

**Decision**: `ts-packages/grpc` builds with `tsup` (centralised at root) rather than per-package `tsc`.

**Why**: matches the working template setup, and produces dual ESM+CJS output without each package owning its own build config.
