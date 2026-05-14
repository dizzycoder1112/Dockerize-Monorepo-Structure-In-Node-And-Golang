---
plan:                    # optional — if this belongs to a plan, the plan's filename slug (e.g. ts-gql-server-migration)
phase:                   # optional — phase number inside the plan (e.g. 3)
depends_on: []           # optional — slugs of tickets that must finish first (e.g. ["static-verification"])
status: open             # open / in-progress / blocked
created: YYYY-MM-DD
updated: YYYY-MM-DD
---

# {Ticket title}

## Context
Why are we doing this? Background, motivation, upstream decisions, related PR / Linear ticket.
If it belongs to a plan, sketch the relationship here (the plan doc has the full context; this is just a pointer).

## Acceptance criteria
- [ ] criterion 1
- [ ] criterion 2
- [ ] criterion 3

## Out of scope
Explicit exclusions to prevent scope creep.

## Notes
In-flight notes, decisions, blockers, links to other tickets.
Append new content at the bottom; do not delete prior entries (preserve the trail).

---

## Ticket conventions

### When to open a ticket
One focused unit of work whose status is worth tracking independently. Implementation (a PR), verification (manual E2E / running a command), or research (a spike) all qualify. If the work is small enough that tracking is pointless (a < 30 min fix), you don't have to open one.

### When NOT to open a ticket — use a plan instead
If the work is large enough to split into ≥ 2 independent tickets with ordering / dependencies → open a plan with phases. A phase only makes sense with ≥ 2 tickets (otherwise just write a single ticket). See `docs/plans/_template.md`.

### Filename & ID
- Filename: `NNN-slug.md`; slug is a descriptive kebab-case phrase (2–4 words).
- **No frontmatter `id:` field** — the number is the filename prefix, the slug is the stable identifier.
- **Numbering rule**: new ticket number = `max(existing active + backlog numbers) + 1`. **Always forward-append, never backward-rename.**
- **When numbers reset**: once all plans / tickets are cleared out (context fully drained), numbering naturally restarts from low values; don't deliberately renumber.
- **Slug is immutable**: once committed, do not rename the slug (it breaks `depends_on` and plan doc links).

### plan / phase linkage
Tickets belonging to a plan should set `plan:` + `phase:` in frontmatter, so `grep "phase: 3" docs/tickets/{active,backlog}/*.md` finds all tickets in a phase.

### Dependencies (`depends_on`)
- List the **slugs** of tickets that must finish first (a deleted ticket counts as done).
- When agents or humans read the queue, **a ticket is not runnable while any `depends_on` slug is still present**, even if its status is `open`.
- No need to manually sync state — the slug doesn't change, and "is the dep done?" is answered by "does that ticket file still exist?".

### Location
- `docs/tickets/active/`: runnable now or in progress (`depends_on` all done; status is open / in-progress / blocked).
- `docs/tickets/backlog/`: waiting on deps or pushed back (`depends_on` not yet satisfied).
- There is no `archive/`.

### Status enum
- `open`: in active/, waiting to be picked up.
- `in-progress`: in active/, being worked on.
- `blocked`: in active/, waiting on an external dependency.

### State transitions
- `backlog/` → `active/`: `mv` once deps are satisfied.
- `open` ↔ `in-progress` ↔ `blocked`: edit frontmatter only, don't move the file.
- **Done = delete the file**: tick the corresponding phase checkbox in the plan doc, then delete. History lives in git log and the plan doc.
- Abandon: note it in the plan doc's Notes section, then delete.
- **Avoid mixing "move + heavy edits" in the same commit** (it breaks git rename detection).
