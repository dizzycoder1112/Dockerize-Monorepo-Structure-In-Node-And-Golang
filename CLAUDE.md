# Claude Notes

Open-source monorepo template showcasing Dockerized Node + Go services across multiple backend architectures.

See `docs/architecture.md` for app inventory, design conventions, runtime behaviour, and key decisions (ADRs).

## In-flight work

Work is tracked across three containers — **plan** / **ticket** / **issue** — under `docs/`.

```
docs/
├── plans/        # multi-ticket strategic container (active/ + archive/)
├── tickets/      # one focused unit of work = one PR / verification / spike (active/ + backlog/; done = delete file)
└── issues/       # observed problems (flat layout, delete when handled)
```

See each directory's `_template.md` for full conventions and frontmatter format.

### Concepts

- **Ticket**: one focused unit of work, the smallest tracking unit. Filename `NNN-slug.md`; the slug is the stable identifier.
- **Plan**: long-running work spanning multiple tickets. Filename is the slug, no number.
- **Phase**: a `## Phase N` section inside a plan (not its own file); only meaningful with ≥ 2 tickets.
- **Issue**: an observed problem (observed ≠ decided to fix). Can be promoted to a ticket or deleted.

### Decision flow

```
new task
├── solvable in one PR?           → open a ticket
├── multiple PRs with ordering?   → open a plan, slice into phases
└── multiple PRs, no phase shape? → open a plan, list tickets flat
```

### Ticket numbering & dependencies

- New number = `max(existing active + backlog numbers) + 1` — **always forward-append, never backward-rename**.
- A completed ticket is deleted → its slug no longer existing in `depends_on` means that dep is done.
- A ticket whose `depends_on` still includes any missing slug is considered **not runnable**, even if it sits in `active/`.

### Agent behaviour

When the user:

- mentions a ticket number or slug ("that #003", "the reserved-subdomains one") → read `docs/tickets/{active,backlog}/`
- mentions a plan slug, or a topic matching a plan under `docs/plans/active/` → read that plan
- mentions an "issue", "bug", or "that problem" matching something under `docs/issues/` → read that issue
- raises a topic matching an active ticket / plan / issue → same — read first, then act

While working on a ticket or plan, if a new out-of-scope problem turns up → open an **issue**, not a ticket (unless you've already decided to fix it now — then it's a ticket).

New tickets and plans are created from the corresponding `_template.md`.
