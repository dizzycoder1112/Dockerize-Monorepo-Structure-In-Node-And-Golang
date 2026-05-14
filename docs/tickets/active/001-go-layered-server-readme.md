---
status: open
created: 2026-05-13
updated: 2026-05-13
---

# Rewrite `go-layered-server` README

## Context

`apps/go-layered-server` still ships the original `go-gin-server` template README ("Go Counter Server") — wrong content after the 3-layer refactor. The README was deliberately deferred during the pre-interview push to keep risk down; picking it up now is safe and is a one-PR job.

## Acceptance criteria

- [ ] README accurately describes the 3-layer architecture (handler → service → repository).
- [ ] Documents how to run zero-infra (`go run`) vs Postgres mode (`DATABASE_URL` set), matching the runtime behaviour described in `docs/architecture.md`.
- [ ] Links out to `docs/architecture.md` for cross-app conventions instead of restating them.

## Out of scope

- Code changes — docs only.
- Adding code examples beyond a minimal "how to run" snippet.

## Notes
