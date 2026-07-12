# AIRFA Website 2.0 — agent instructions

**Before doing anything: read `agent/START-HERE.md`, then `agent/STATUS.md`.** They define the onboarding order, the working rules, and the current state of the project. All specs live in `agent/`.

Non-negotiables (details in START-HERE.md):

- The design mockups in `agent/design/` are the source of truth for all UI — follow them strictly; verify visually (screenshot vs mockup) before calling a page done.
- Backend is **Go** (chi + pgx + sqlc). Frontend is **Next.js + TypeScript + Tailwind**. No Prisma, no ORMs.
- Update `agent/STATUS.md` at the end of every session and tick `agent/roadmap.md` items as they complete. Log decisions in STATUS.md → Decisions.
- The previous 2.0 attempt is void — never resurrect its decisions. `../airfa-website/` (current live site) is a content reference only.
- Site content is Portuguese (pt-PT); code and docs are English; CMS admin UI is Portuguese.
- Commits: local checkpoints, plain messages, never co-author lines.
