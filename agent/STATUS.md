# STATUS — living handover document

> Update this file at the end of every working session. Newest information at the top of each section.

## Current state (2026-07-13)

**Phase: 0 scaffolding complete → starting Phase 1 · Stage A.**

> **⚑ Active working checklist lives in `current_phase.md` at the repo root.** It has the detailed, ticked Step-0 items and the full Stage A steps (schema SQL, sqlc queries, endpoints) with copy-paste commands. Read it first.

The monorepo skeleton is up and verified. Layout: `apps/api` (Go), `apps/web` (Next.js), `packages/shared` (shared TS); root `Makefile` runs everything; `go.work` + npm workspaces.

Done (all green — `go vet`/`go build`/`go test`/`tsc --noEmit`/`next build`):
- `airfa_dev` recreated clean; Go API (`apps/api`) serves `GET /healthz`; `go.work` set.
- Next.js app (`apps/web`, Next 16 + Tailwind v4) with `/api/*` → Go API rewrite.
- `packages/shared` stub linked into web via npm workspaces; `overrides.postcss ^8.5.10` in root `package.json` (keep it).
- Root `Makefile`, `.editorconfig`, `.gitignore`, `.env.example` (both apps); `migrate` + `sqlc v1.31.1` installed; `apps/api/sqlc.yaml` present.

Not yet done (intentionally folded into Stage A): first migration (users/sessions), sqlc generating, and the web→API round-trip page. Nothing is committed to git yet (all new files untracked).

**Working model:** the owner does not code — an agent implements; the **owner runs all commands himself** (agent gives specs/steps + verifies, does not run app builds/servers). Give exact, copy-paste steps.

**Machine facts:** dev DB `postgres://airfa:airfa@localhost:5432/airfa_dev?sslmode=disable`; `~/go/bin` on PATH via `~/.zshrc` (new shells); Go 1.26.4, Node 26.4, Postgres 18.

## Next steps

1. **Phase 1 · Stage A** (see `current_phase.md`): A.1 migrations (users/sessions + content + collections), A.2 sqlc queries + generate, A.3 public read endpoints (`/v1/site`, `/v1/pages/:slug`, `/v1/collections/*`). This also closes the leftover Phase 0 items.
2. **Stage B:** design tokens — Inter via `next/font`, Tailwind colors sampled from the mockups.
3. Then Stages C–H (chrome, blocks, benchmark pages, remaining pages, seed, polish) per `current_phase.md`'s Phase 1 map.

## Open questions for the owner

- [ ] **Hosting**: assumed Docker via Dokploy on your VPS. Confirm.
- [ ] **CMS roles**: proposed Admin + Editor (requirements.md §4). Confirm or simplify to a single admin role for v1.
- [ ] **Versioning depth**: proposed simple published/draft + version history with restore (requirements.md §5). Confirm.
- [x] **Blog & events** — *resolved 2026-07-12:* **both are real content** the association will maintain. Build and seed both the Events/Cartazes and Blog collections fully in Phase 1.
- [x] **Donations (DOAR button)** — *resolved 2026-07-12:* links to an **on-site bank details / IBAN page** (how-to-donate; NIB + transfer info). No external platform for v1.
- [x] **Exact fonts** — *resolved 2026-07-12:* **Inter** (Google Fonts) for all text, self-hosted via `next/font` at build time (no runtime CDN). Verify the giant display numeral ("1895") against mockup-02 during Stage B.

## Decisions

- **2026-07-17 — All non-user-visible strings are English, including API error messages.** Owner's call, reversing the "pt for CMS-facing messages" wording in architecture.md §2.1 (§2.1 updated). Rule: **only text that literally renders on the public page or in the CMS admin UI is Portuguese.** Logs, Go errors, and the API error envelope's `message`/`fields` are English — they're developer-facing. The frontend never shows an API `message` raw; it switches on the stable `code` and renders its own pt-PT copy, which is also where any translation happens. **Corollary: `code` is the API contract** — keep codes stable and machine-readable, and never have the frontend match on `message` text. Applied across `apps/api/internal/http/*` (5 strings).

- **2026-07-12 — Fresh start; disregard previous 2.0 iteration.** Owner explicitly instructed that all prior decisions/artifacts from the earlier attempt are void. `../airfa-website/` (current live site code) remains the content reference.
- **2026-07-12 — Stack: Go API + Next.js frontend.** Owner wants Go on the backend; Next.js chosen over a plain React SPA for SSR/SEO/link previews on a public association site. Frontend is React + TypeScript per owner preference.
- **2026-07-12 — Portuguese-only content** for launch; schema should not block adding locales later (assumed by agent, silently confirmable).
- **2026-07-12 — Fonts: Inter everywhere.** Owner confirmed the design uses Inter (Google Fonts). Self-host via `next/font/google` (build-time, no runtime CDN — GDPR). Replaces the Archivo/Playfair proposal. Display numerals ("1895") also Inter unless measuring mockup-02 shows an obvious serif — flag to owner if so.
- **2026-07-12 — Blog & Events are real content types.** Owner confirmed the association will maintain both. Both collections (posts, events) are fully built and seeded in Phase 1, not placeholders.
- **2026-07-12 — DOAR button → on-site IBAN page.** The donate button links to a page on the site with bank/NIB details and how-to-donate text. No external donation platform in v1 (DOAR remains a link, not a checkout — consistent with requirements §7).
- **2026-07-12 — Monorepo with `apps/` + `packages/` layout.** Owner wants an intentional monorepo. Structure: `apps/api` (Go), `apps/web` (Next.js), `packages/shared` (shared TS block types + fixtures). Root `Makefile` is the unified task runner across both languages; npm workspaces manage the JS side; a root `go.work` manages Go. Go module path becomes `github.com/RMLoureiro/airfa-website-2.0/apps/api`. Turborepo deferred until multiple JS build targets exist. Repo layout recorded in architecture.md §1.1.
- **2026-07-12 — Main nav label: "Actividades".** Resolves the mockup discrepancy ("Actividades" vs "O que fazemos"); matches the current site and content-inventory naming.

## Blockers

None.

## Session log

- **2026-07-13** — Built the monorepo foundation (owner ran all commands; agent gave steps + verified). Created `apps/api` (Go/chi, `/healthz`), `apps/web` (Next 16/TS/Tailwind v4 + `/api/*` rewrite), `packages/shared`, root npm workspace + `go.work` + `Makefile` + `.editorconfig`/`.gitignore`/`.env.example`. Installed `migrate` + `sqlc`; added `sqlc.yaml`. Fixed a chi route-slash panic; fixed a Next-bundled postcss advisory via root `overrides` (clean reinstall → 0 vulns, `next build` passes); added `~/go/bin` to `~/.zshrc`. All checks green. Wrote `current_phase.md` (root) as the active working checklist through Stage A. Left first migration / sqlc-generate / round-trip page for Stage A; nothing committed to git yet.
- **2026-07-12 (d)** — Owner walkthrough of Phase 0 + Phase 1 (guidance only, no app code). Dropped and recreated `airfa_dev` clean (removed the discarded Prisma tables). Resolved 4 open questions with the owner: fonts = Inter (self-host via next/font), Blog & Events = both real, DOAR = on-site IBAN page, nav label = "Actividades". Updated STATUS (decisions + open questions), design-spec (typography), requirements (blog scope, DOAR, nav label). Adopted a monorepo (`apps/` + `packages/`, npm workspaces + `go.work`, Make as unified runner) — architecture.md §1.1 + layout refs updated.
- **2026-07-12 (c)** — Consistency pass: fixed requirements.md § references in Open questions; specified homepage slug convention (`inicio` → `/`), image variant sizes (320/960/1920), font self-hosting via next/font, admin-UI-has-no-mockups rule; verified toolchain (Go 1.26.4, Node 26.4).
- **2026-07-12 (b)** — Detail pass: added API conventions, env vars, and canonical block-tree JSON to architecture.md; layout metrics to design-spec.md; living-document rules to all agent docs; CLAUDE.md created; project plugins enabled (frontend-design, playwright, chrome-devtools-mcp).
- **2026-07-12** — Project restarted from scratch. Wrote documentation set, saved design assets, initialized planning. No code yet.
