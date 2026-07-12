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

## Tooling & skills (use them — they're set up for this project)

Project-scope plugins are enabled in `.claude/settings.json` (they load at session start; if you enabled them mid-session, restart the session):

- **frontend-design** — invoke when building any public-site UI. The mockups must be matched *strictly*; this skill raises the design quality bar. Use it before writing page/component code.
- **playwright** — browser automation. Use it to (a) visually verify built pages against the mockups (screenshot → compare side by side), and (b) write the e2e round-trip tests required by the roadmap (login → edit → publish → verify).
- **chrome-devtools-mcp** — inspect the running app (console errors, network, DOM) when debugging. Note: no Docker locally; a bundled Chrome exists at `../airfa-website/chrome/linux-150.0.7871.24/chrome-linux64/chrome` if a browser binary is needed.
- ~~prisma~~ — deliberately disabled; the backend is Go + sqlc, **not** Prisma. Don't re-enable it.

Built-in skills that are expected practice here:

- **/verify** — before committing any nontrivial change, exercise the affected flow end-to-end (run the app, click through), not just typecheck.
- **/code-review** — run after completing a feature-sized chunk; fix confirmed findings before the checkpoint commit.
- **/simplify** — run occasionally after large additions to keep the codebase lean.

**Visual verification is mandatory for Phase 1 pages:** render the page, screenshot it at desktop width, and compare against the corresponding mockup in `design/` before calling it done. Record the comparison result in STATUS.md.

## Keeping these documents alive

These files are the project's memory between agents. Stale docs are worse than no docs. Rules:

| File | Update when… |
|---|---|
| `STATUS.md` | **Every session, no exceptions** — done / half-done / next / decisions / blockers. |
| `roadmap.md` | An item completes (tick it) or scope changes (edit items + note in STATUS.md → Decisions). |
| `requirements.md` | The owner confirms an open question or changes scope. Replace *(proposal — confirm)* markers with the decision + date. |
| `architecture.md` | Any technical decision deviates from what's written. Update the doc **and** log why in STATUS.md → Decisions. |
| `design-spec.md` | Exact tokens are sampled from mockups (replace ≈ values), a new block type is added (extend §6 table), or the owner rules on a design question. |
| `content-inventory.md` | Content gets migrated (tick/annotate), or facts are verified/corrected during seeding. |
| `START-HERE.md` | Tooling, rules, or read-order change. |

If you create a new significant doc, link it from this file and from README.md.

## Quick facts

- Backend: Go (chi, pgx, sqlc, golang-migrate) — see architecture.md
- Frontend: Next.js App Router + TypeScript + Tailwind — public site (SSR) + `/area-reservada` admin
- DB: PostgreSQL (local dev instance available on this machine, port 5432)
- Deploy target: Docker via Dokploy on the owner's VPS (Phase 3)
- Content reference: `../airfa-website/` (static build of current airfa.pt) and https://www.airfa.pt
