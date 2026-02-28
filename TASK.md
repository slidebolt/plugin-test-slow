# TASK

## Scope
Harden runtime behavior and test hygiene for this repository.

## Constraints
- No git commits or tags from subprocesses unless explicitly requested.
- Keep changes minimal, testable, and production-safe.
- Prefer deterministic shutdown/startup behavior.

## Required Output
- Small PR-sized patch.
- Repro steps.
- Validation commands and expected results.
- Known risks/limits.

## Priority Tasks
1. Preserve deterministic behavior for harness stress validation.
2. Ensure clean exit under supervisor stop.
3. Keep failures intentional and observable (for flaky/slow semantics).

## Done Criteria
- Test plugins accurately validate supervisor hardening without orphan processes.

## Validation Checklist
- [ ] Build succeeds for this repo.
- [ ] Local targeted tests (if present) pass.
- [ ] No new background orphan processes remain.
- [ ] Logs clearly show failure causes.
