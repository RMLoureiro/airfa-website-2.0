# STATUS — living handover document

> Update this file at the end of every working session. Newest information at the top of each section.

## Current state (2026-07-12)

**Phase: 0 — not started (documentation only).**

The repository contains only documentation and design assets. No code has been written. The previous 2.0 attempt was discarded entirely (owner's instruction); this doc set is a fresh start.

Done so far:
- Design mockups + logo saved to `agent/design/`.
- Full documentation set written (requirements, architecture, design spec, content inventory, roadmap, README).
- Stack decided with the owner: **Go API + Next.js frontend + Postgres** (see Decisions).

## Next steps

1. Owner reviews the documentation set and answers the open questions below.
2. Begin **Phase 0** (roadmap.md): repo scaffolding — `api/` Go module, `web/` Next.js app, migrations tooling, dev database, `make dev` workflow.
3. Then Phase 1: content schema + public site with the new design.

## Open questions for the owner

- [ ] **Hosting**: assumed Docker via Dokploy on your VPS. Confirm.
- [ ] **CMS roles**: proposed Admin + Editor (requirements.md §5). Confirm or simplify to a single admin role for v1.
- [ ] **Versioning depth**: proposed simple published/draft + version history with restore (requirements.md §6). Confirm.
- [ ] **Blog & events**: the homepage mockup shows a Blog section and event "cartazes". Are these real content types the association will feed, or placeholders? (Affects Phase 1 scope.)
- [ ] **Donations (DOAR button)**: where should it link? External platform, IBAN info page, or something else?
- [ ] **Exact fonts**: design-spec.md proposes families matching the mockups; confirm or provide the original design's font names if known.

## Decisions

- **2026-07-12 — Fresh start; disregard previous 2.0 iteration.** Owner explicitly instructed that all prior decisions/artifacts from the earlier attempt are void. `../airfa-website/` (current live site code) remains the content reference.
- **2026-07-12 — Stack: Go API + Next.js frontend.** Owner wants Go on the backend; Next.js chosen over a plain React SPA for SSR/SEO/link previews on a public association site. Frontend is React + TypeScript per owner preference.
- **2026-07-12 — Portuguese-only content** for launch; schema should not block adding locales later (assumed by agent, silently confirmable).

## Blockers

None.

## Session log

- **2026-07-12** — Project restarted from scratch. Wrote documentation set, saved design assets, initialized planning. No code yet.
