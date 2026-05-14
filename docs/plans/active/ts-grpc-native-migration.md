---
status: draft
created: 2026-05-13
updated: 2026-05-13
---

# Migrate TS side off connect-es to native gRPC

## Goal & Context

Today: the TS server speaks Connect protocol via `connectNodeAdapter`, while the Go server speaks native gRPC. The shared `createGrpcTransport` client masks the difference, but the TS server protocol is still a leaky abstraction — the underlying wire isn't uniform across servers.

Target: native gRPC end-to-end via `@grpc/grpc-js` + `ts-proto` on the TS side. Same wire as the Go side, same debugging tooling everywhere (grpcurl / reflection / wireshark).

End state criteria:
- `@bufbuild/protobuf`, `@connectrpc/connect`, `@connectrpc/connect-node`, `protoc-gen-connect-es`, `protoc-gen-es` are all gone from the dependency graph.
- `pnpm run build` green across the monorepo.
- End-to-end curl through `ts-restful-api` → `ts-grpc-demo` still returns `{"message":"You said X"}`.
- The single `createGreeterClient` factory still hits both `ts-grpc-demo` (`:50051`) and `go-grpc-demo` (`:50052`) with identical API.

## Decisions

- Q: stay on connect-es (pinned to v1) vs migrate to native gRPC / A: migrate — connect-es v1 is the last release on that path; the v2 direction is `protoc-gen-es` schema-only without a separate connect plugin. Staying on v1 indefinitely is a dead end.
- Q: `ts-proto` vs other TS gRPC codegens / A: `ts-proto` — it outputs idiomatic interfaces (versus class-based messages) and is the conventional pairing with `@grpc/grpc-js`.
- Q: migrate the Go side too / A: no — Go side already uses native gRPC (`protoc-gen-go` + `protoc-gen-go-grpc`).

## Risks & Open questions

- `ts-proto`'s generated types diverge from `@bufbuild/protobuf`'s. The public surface of `ts-packages/grpc` changes shape, so `apps/ts-restful-api` and `apps/ts-grpc-demo` need coordinated rewrites in the same PR.
- Anything currently relying on Connect-over-HTTP fallback at the server would break after migration. Today no consumer does — confirm before merging.
- Estimated 1–2 h of breaking changes once started.

## Out of scope

- Go side rewrite (already native gRPC).
- Changing the `proto/hello.proto` schema itself.
- Adding new RPC methods.

## Phases

### Phase 0: toolchain + codegen swap [📋 open]
**Scope**: replace bufbuild/connect-es plugins with `ts-proto`, regenerate TS protobuf output, update the `ts-packages/grpc` dependency graph.
**Tickets**:
- [ ] (TBD) — swap plugins in root `package.json` + `proto/buf.gen.yaml`, regenerate via `pnpm run buf:gen`, commit the new generated files
- [ ] (TBD) — update `ts-packages/grpc/package.json`: drop `@bufbuild/protobuf` / `@connectrpc/connect` / `@connectrpc/connect-node`, add `@grpc/grpc-js`, run `pnpm install`

### Phase 1: rewrite TS consumers against `@grpc/grpc-js` [📋 open]
**Scope**: rewrite the shared grpc wrapper, the demo handler, and the REST API's client usage against the new generated types.
**Tickets**:
- [ ] (TBD) — rewrite `ts-packages/grpc/src/{clientFactory,serverFactory,index}.ts` for `@grpc/grpc-js`
- [ ] (TBD) — rewrite `apps/ts-grpc-demo/src/handlers/sayHello.handler.ts` (promise/callback style), update `apps/ts-restful-api/src/services/index.ts` + `repositories/user.repository.ts`, verify end-to-end curl roundtrip + `pnpm run build` green

## Why deferred

Too risky pre-interview (2026-04-30 — see [[interview_2026_04_30]]). Pick up post-interview.
