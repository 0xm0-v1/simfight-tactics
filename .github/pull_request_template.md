# PR Title
<type>[<scope>]: <short summary> (#<issue_id>)

## What
Short summary of the change. Keep it crisp (1–3 sentences).

## Why
Link the intent and the issue.
Closes #<issue_id>

## How
Key implementation points:
- Major decisions / data structures
- Backward-compat notes
- Flags/config toggles (if any)

## Tests
- [ ] Unit / integration scenarios listed
- [ ] Deterministic with seed (if RNG used)
- [ ] Added/updated fixtures & golden files (if applicable)

## Observability
- [ ] Structured logs added/updated
- [ ] Metrics (mean/median/std or relevant) covered
- [ ] Docs: mechanics/README updated if behavior changed

## UI (skip if N/A)
- [ ] SSR/HTMX rendering checked
- [ ] Screenshots / before-after attached

## Risk & Rollout
- [ ] Breaking changes? If yes, described and gated
- [ ] Feature flagged or reversible
- [ ] Migration steps (if any)

## Checklist (Ready to merge)
- [ ] PR title follows convention: `<type>[<scope>]: <summary> (#<issue>)`
- [ ] Labels applied (enhancement/documentation/…)
- [ ] Linked issue auto-closing keyword present (`Closes #…`)
- [ ] Branch up to date with `main`
- [ ] All checks green (CI/lint/tests)

