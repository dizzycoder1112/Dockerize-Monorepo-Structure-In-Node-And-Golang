---
status: draft       # draft / in-progress / blocked / completed / abandoned
created: YYYY-MM-DD
updated: YYYY-MM-DD
---

# {Plan title}

## Goal & Context
Why is this plan starting? Background, upstream decisions, expected end state.

## Decisions
Cross-phase design decisions. One line each: "Q: xxx vs yyy / A: chose yyy because ...".

## Risks & Open questions
- Known risks
- Open design questions

## Out of scope
Explicit exclusions to prevent scope creep.

## Phases

### Phase 0: {phase name} [📋 open]
**Scope**: what this phase does.
**Tickets**:
- [ ] #NNN — ticket description

### Phase 1: {phase name} [📋 open]
**Scope**: ...
**Tickets**:
- [ ] #NNN — ...

---

## Plan conventions

- **Location**: two states — active / archive. Draft also lives in `active/` (draft is not abandoned — it's just not started yet).
  - `docs/plans/active/`: status = draft / in-progress / blocked
  - `docs/plans/archive/`: status = completed / abandoned
- **Filename**: slug-based (`ts-gql-server-migration.md`), no numbering.
- **Status semantics**:
  - `draft`: scoped but not started. Typical reasons: waiting on prerequisites / still deciding whether to do it / not yet scheduled.
  - `in-progress`: actively being worked on.
  - `blocked`: in flight but stuck waiting on something external.
  - `completed`: finished.
  - `abandoned`: dropped.
- **Phase count**: a phase only makes sense with ≥ 2 tickets. Merge or dissolve any phase with only one ticket.
- **Phase status emoji**: `📋 open` / `🚧 in-progress` / `🚫 blocked` / `✅ done`
- **State transitions**:
  - draft → in-progress: update frontmatter when work starts.
  - in-progress ↔ blocked: just update frontmatter and the phase emoji.
  - Closing out: move to `archive/`, set frontmatter `status: completed`, and add `closed: YYYY-MM-DD`.
  - **Avoid mixing "move + heavy edits" in the same commit** (it breaks git rename detection).
