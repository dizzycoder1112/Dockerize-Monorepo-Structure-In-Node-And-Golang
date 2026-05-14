---
area:                    # impact area: which app / package (e.g. ts-game / ts-packages/grpc)
severity: low            # low / medium / high
created: YYYY-MM-DD
---

# {Issue title}

## Observation
What was seen.

## Reproduction
How to trigger it / where it shows up.

## Suspected cause
A guess at the root cause (no need to verify).

## Disposition
Undecided / will fix / won't fix.

---

## Issue conventions

### When to open an issue
You observed a problem, bug, tech debt, or annoyance — and either you're not sure whether to fix it, or you've decided not to fix it now but want a record.

### Difference from a ticket
- **Ticket**: work you've decided to do.
- **Issue**: a problem you observed (may or may not get fixed later).
- An issue can be promoted to a ticket (when you decide to fix it).

### Filename
`<area>-<short-slug>.md`, e.g. `ts-game-typecheck-vitest-vite-type-def.md`. **No number** (volume is low, no ordering needed, no clash with ticket numbers).

### Location
Flat under `docs/issues/`, **no active / closed subdirectories**. Delete the file when handled.

### State transitions
- **Promote to ticket**: open a new ticket whose Context cites `originated from docs/issues/<file>.md`, then delete the issue file.
- **Won't fix**: note the won't-fix reason in the commit message, then delete the issue file.
- **Undecided**: leave it.

### Not archived
History lives in git log + commit messages.
