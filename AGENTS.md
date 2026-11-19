# AGENTS.md — Engineering Mastery project policy

Work directly in the current session on this original project. Keep the Engineering Mastery manifest and progress report current when working in that workspace.

## Scope and ethics
- Never access, rewrite, delete or search outside the active project directory.
- Reference products are inspiration, not source-copy instructions. Respect licenses, IP and trademarks.
- No employer-owned source, production secrets, private personal data or real customer records.
- Security testing must use only local authorized fixtures. No external target testing.
- Do not provision paid cloud services, purchase assets, publish publicly, or rotate credentials.
- Do not push or merge to GitHub from the IDE. A separately authorized private GitHub workflow or the human handles remote operations.

## Engineering behavior
- Inspect existing code, tests, current git status and decisions first.
- Maintain README.md, docs/TASKS.md, docs/EVIDENCE.md and PROJECT_STATE.json.
- Complete one small, verifiable milestone at a time and continue while the IDE session and daily budget permit.
- Produce real runnable code; use reproducible synthetic fixtures and test first where sensible.
- Every nontrivial public API needs contract and misuse tests. Document decisions with ADRs.
- Create scripts/verify.ps1 with REAL nontrivial checks. Never make it an always-success shim.
- Run all project checks you can. Mark a test as not run if missing tools; do not claim passing.
- Commit only actual completed changes, with one logical purpose per commit. No empty/vanity commits.
- Record test output and any honest limitations in docs/EVIDENCE.md.
- If blocked, record why and what exact human action is necessary.
- Set `PROJECT_STATE.json` status to `release_review` after the working MVP, required verification, and release materials are complete. Wait for explicit user acceptance before recording completion or advancing the queue.
- End each invocation with a short summary: modified files, commits, tests, next task, blockers.
