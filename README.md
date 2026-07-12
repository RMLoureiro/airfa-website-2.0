# AIRFA Website 2.0

Rebuild of the **Academia de Instrução e Recreio Familiar Almadense (AIRFA)** website — [www.airfa.pt](https://www.airfa.pt) — with a modern design and a full self-service CMS so non-technical members can edit every part of the site (pages, menus, blocks, media) without touching code.

> **New to this project (human or AI agent)?** Read **[`agent/START-HERE.md`](agent/START-HERE.md)** first. It explains the current state, the documentation map, and the working rules.

## What this is

- **Public site** — all content of the current airfa.pt, re-implemented with the new design (mockups in [`agent/design/`](agent/design/)). Portuguese-only for now.
- **CMS ("Área Reservada")** — login-protected admin where editors can add/remove/reorder blocks, create pages and menu entries, link pages to each other, upload media, and restore previous versions. Think "WordPress-like", but purpose-built and typed.

## Stack

| Layer | Choice |
|---|---|
| Backend | **Go** (chi router, pgx + sqlc, golang-migrate) — REST API, auth, content model |
| Frontend | **Next.js** (App Router, React + TypeScript, Tailwind CSS) — public site (SSR) + admin UI |
| Database | **PostgreSQL** |
| Deploy | Docker containers via **Dokploy** on a personal VPS |

Why this split: the public site needs SEO and link previews (SSR via Next.js), while Go owns all data, business logic, and auth. Full rationale in [`agent/architecture.md`](agent/architecture.md).

## Repository layout (planned)

```
api/            Go backend (REST API, migrations, media storage)
web/            Next.js frontend (public site + /area-reservada admin)
agent/          Project documentation for humans and AI agents
  START-HERE.md   Onboarding — read this first
  STATUS.md       Living status: what's done, what's next (keep updated!)
  requirements.md Functional requirements
  architecture.md System design, data model, API surface
  design-spec.md  Design system + per-mockup component spec
  content-inventory.md  Pages/content to migrate from airfa.pt
  roadmap.md      Phased implementation plan with checkboxes
  design/         Mockups + logo (source of truth for the visual design)
```

## Tooling

Project-scope Claude Code plugins are enabled in `.claude/settings.json`: **frontend-design** (UI implementation quality), **playwright** (visual verification + e2e), **chrome-devtools-mcp** (runtime debugging). `CLAUDE.md` routes agents to the onboarding doc. Expected practices (verify before commit, code review per feature) are defined in `agent/START-HERE.md`.

## Getting started (once code exists)

Phase 0 of the roadmap creates the scaffolding; until then there is nothing to run. The dev environment assumes a local Postgres (see `agent/architecture.md` → Development environment).

## Reference material

- Old static site source: `../airfa-website/` (this machine) and https://www.airfa.pt
- Design mockups: `agent/design/mockup-*.{jpeg,png}`
- Logo/crest: `agent/design/logo-airfa.png`
