# 🧱 Dockerize Monorepo Structure with Golang and TypeScript

A monorepo template for high-performance services that pairs Go and Node side-by-side: shared `proto/` definitions, generated clients on both sides, layered + DDD Go server examples, and an Express TypeScript app — all wrapped in Docker.

> ✨ Built for microservices, modular APIs, and teams scaling shared utilities across language boundaries.

## ⚡️ Key Features

1. **Containerized DX** — Dockerized monorepo keeps local environment clean and consistent with production.
2. **Hot reload & build orchestration** — turbo watches workspace changes and rebuilds dependent packages.
3. **Two Go backend architectures**, side by side:
   - `apps/go-layered-server` — handler / service / repository (3-layer)
   - `apps/go-ddd-server` — domain / app / infra / interfaces (DDD)

   Both ship with a composition root in `internal/factory/`, sentinel errors, and an in-memory repository so they boot zero-config.
4. **Type-safe gRPC pipeline** — single `proto/*.proto` source, `pnpm run buf:gen` produces matching Go (`go-packages/grpc`) and TypeScript (`ts-packages/grpc`) clients. Native gRPC interop locked to HTTP/2.
5. **Shared Node packages** — `logger`, `db`, `grpc`, `shared`, `rabbitMQ` consumed by every TS app via pnpm workspace.

## 📂 Project Structure

```
root/
├── apps/
│   ├── go-ddd-server/         DDD Go server (in-memory + optional Postgres)
│   ├── go-layered-server/     3-layer Go server (in-memory + optional Postgres)
│   ├── ts-grpc-demo/          ConnectRPC gRPC server (HTTP/2)
│   └── ts-restful-api/        Express REST API
├── go-packages/
│   ├── grpc/                  Generated Go clients + helpers
│   ├── logger/                Unified logger interface (Console / Zap / Sentry strategies)
│   └── rabbitMQ/              Producer / consumer factories
├── ts-packages/
│   ├── db/                    pg + kysely
│   ├── grpc/                  Generated TS clients (Connect / gRPC transports, H2-only)
│   ├── logger/                Pino-flavoured logger
│   ├── rabbitMQ/              amqplib wrapper
│   └── shared/                cross-app constants & utils
├── proto/
│   ├── buf.yaml
│   ├── buf.gen.yaml           emits Go + TS in one shot
│   ├── hello.proto
│   └── eliza.proto
├── kubernetes/
├── docker-compose.local.yml
├── package.json
├── pnpm-workspace.yaml
└── turbo.json
```

## 🛠 Usage

### 1. Install

```bash
pnpm install
```

### 2. Start dev (Docker)

> Requires Docker running locally. Watches workspace changes and reloads via turbo.

```bash
pnpm run start:dev
```

### 3. Build everything

```bash
pnpm run build
```

### 4. Run a Go server directly

```bash
cd apps/go-layered-server && go run ./cmd/main.go
# or
cd apps/go-ddd-server   && go run ./cmd/main.go
```

Both default to in-memory repositories. Set `DATABASE_URL` to flip the layered server onto its pgx-backed Postgres path.

## 🔌 gRPC code generation

`buf` is bundled via npm — no Homebrew step required. The pipeline generates Go + TS in one command:

```bash
pnpm run buf:gen
```

- Edit `proto/*.proto`
- Re-run `pnpm run buf:gen`
- TS clients land in `ts-packages/grpc/src/proto/`, Go clients in `go-packages/grpc/pb/<package>/`

For the Go side you need `protoc-gen-go` and `protoc-gen-go-grpc` once globally:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 💻 Contribution

Fork, improve, send a PR. Goal: make scalable Go + Node monorepos boring to start.
