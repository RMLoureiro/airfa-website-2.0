# START HERE — Agent Onboarding

You are working on **AIRFA Website 2.0**: a rebuild of www.airfa.pt with a modern design and a self-service CMS. The owner (Ricardo) directs the work; agents implement it.

## Read in this order

1. **`STATUS.md`** — where the project is right now, what's in progress, what's next. *Always start here.*
2. **`roadmap.md`** — the phased plan with checkboxes.
3. **`requirements.md`** — what we're building (functional requirements, CMS behavior).
4. **`architecture.md`** — how we're building it (stack, data model, API, repo layout).
5. **`design-spec.md`** — the visual design system and per-page component specs. The mockups in `design/` are the source of truth; follow them **strictly**.
6. **`content-inventory.md`** — the airfa.pt content that must exist at launch.

## Hard rules

- **The design mockups win.** `design/mockup-01-homepage.jpeg`, `design/mockup-02-banda.png`, `design/mockup-03-estatutos.png`. When code and mockup disagree, the mockup is right. Open questions about the design go to the owner, not to your imagination.
- **Keep `STATUS.md` updated.** Every working session ends by updating it: what was done, what's half-done, exact next steps, any new decisions or blockers. The next agent may have zero context beyond these files.
- **Tick off `roadmap.md`** items as they are completed. If scope changes, update the roadmap and note the change in STATUS.md.
- **Decisions log.** Non-trivial technical decisions get a dated entry in STATUS.md → Decisions. Don't silently deviate from architecture.md — update it and log why.
- **Ignore any leftovers of the previous 2.0 attempt.** There was an earlier iteration of this project (different stack). The owner has said to disregard it entirely. This documentation set is the only truth. The old *live-site* codebase at `../airfa-website/` is valid as a **content** reference only — not for design or architecture.
- **Commits:** local commits at natural checkpoints. Never add co-author lines or tool attributions.
- **Owner does not code.** Explain what you did in plain terms; he reviews behavior, not diffs.
- **Language:** site content is Portuguese (pt-PT). Code, comments, and docs are English. CMS admin UI text: Portuguese.

## Quick facts

- Backend: Go (chi, pgx, sqlc, golang-migrate) — see architecture.md
- Frontend: Next.js App Router + TypeScript + Tailwind — public site (SSR) + `/area-reservada` admin
- DB: PostgreSQL (local dev instance available on this machine, port 5432)
- Deploy target: Docker via Dokploy on the owner's VPS (Phase 3)
- Content reference: `../airfa-website/` (static build of current airfa.pt) and https://www.airfa.pt
