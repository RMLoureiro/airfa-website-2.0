# Architecture

## 1. Overview

Two applications + one database, deployed as Docker containers behind Dokploy's reverse proxy:

```
┌────────────┐   HTML/SSR    ┌─────────────┐   REST/JSON   ┌──────────┐
│  Browser   │ ────────────▶ │  web (Next) │ ────────────▶ │ api (Go) │
│            │  /api/* proxied to Go ──────▶│              │──▶ Postgres
└────────────┘               └─────────────┘               └──────────┘
                                                             └──▶ media volume
```

- **`api/` — Go service.** Owns the database, all business logic, auth, media uploads. Exposes a versioned REST API (`/v1/...`). The only writer to Postgres.
- **`web/` — Next.js app.** Renders the public site server-side by fetching published content from the API; hosts the admin SPA-style UI at `/area-reservada` which talks to the API (through a same-origin `/api` rewrite so cookies are first-party).
- **PostgreSQL** — content, users, sessions, versions.
- **Media volume** — uploaded files on disk (`/data/media`), served by the Go API (`/media/...`) with cache headers. Simple, backup-friendly; can move to S3-compatible storage later without schema changes (paths are opaque URLs to the frontend).

## 2. Backend (Go)

- **Go** latest stable; standard project layout:
  ```
  api/
    cmd/server/main.go
    internal/
      http/        handlers, middleware, router (chi)
      auth/        sessions, password hashing (argon2id), RBAC
      content/     pages, blocks, menus, collections, versions
      media/       upload, image variants, static serving
      db/          sqlc-generated code, pgx pool
    migrations/    golang-migrate SQL files
    sqlc.yaml
  ```
- **Libraries** (deliberately boring):
  - Router: `go-chi/chi` (stdlib-compatible, middleware ecosystem)
  - DB: `jackc/pgx/v5` + `sqlc` (typed queries, no ORM magic)
  - Migrations: `golang-migrate/migrate` (SQL files in repo)
  - Auth: `alexedwards/argon2id` for hashing; DB-backed sessions with httpOnly, Secure, SameSite=Lax cookie (revocable — better than JWT for a CMS)
  - Validation: `go-playground/validator` on request DTOs; block payloads validated against per-type Go structs
  - Images: `disintegration/imaging` (or `bimg` if libvips acceptable) for resize variants
- **API surface (v1, illustrative):**
  - Public (no auth): `GET /v1/site` (settings, menus), `GET /v1/pages/:slug` (published block tree), `GET /v1/collections/{events|posts|partners|activities}` with published filter
  - Auth: `POST /v1/auth/login`, `POST /v1/auth/logout`, `GET /v1/auth/me`
  - Admin (session + role): CRUD on pages/blocks (draft state), menus, media, collections, settings, users (Admin only), versions (`GET /v1/pages/:id/versions`, `POST /v1/versions/:id/restore`)
  - `POST /v1/pages/:id/publish` — promotes draft → published, writes a version snapshot, and triggers frontend revalidation (calls a Next.js revalidate webhook with a shared secret)

## 3. Content model

The heart of the system: pages are **trees of typed blocks**, stored as JSONB but validated by Go structs per block type (the same catalog the frontend renders — `design-spec.md` §Blocks).

Core tables (simplified):

```sql
pages(id, slug, title, seo jsonb, parent_id, sort, status, deleted_at, created_at, updated_at)
page_revisions(id, page_id, kind draft|published|snapshot, blocks jsonb, seo jsonb,
               created_by, created_at)          -- current draft & published are rows here
menus(id, zone, items jsonb, updated_at)         -- zone: main | secondary | utility | footer
settings(id, data jsonb)                         -- contacts, socials, doar_url, logo, ...
media(id, kind image|pdf, path, filename, size, width, height, alt, created_at)
events(id, title, poster_media_id, starts_at, ends_at, body jsonb, status, ...)
posts(id, title, slug, cover_media_id, excerpt, blocks jsonb, published_at, status, ...)
partners(id, name, logo_media_id, url, sort)
activities(id, name, category, image_media_id, info jsonb, sort, status)
users(id, name, email, password_hash, role admin|editor, created_at, disabled_at)
sessions(id, user_id, token_hash, expires_at, created_at)
```

Notes:
- **Draft/published** = two revision pointers per page; publishing copies draft → published and appends a `snapshot` revision (history). Restore = copy snapshot → new draft.
- Blocks carry stable `id`s (uuid) inside the JSON so the builder can reorder/patch without diffs.
- Internal links inside block fields are stored as `{ "type": "page", "pageId": ... }` and resolved to slugs at render time — renaming a slug never breaks links.
- A future `locale` column on `page_revisions`/collections is the multi-language path; don't build it now, don't block it.

## 4. Frontend (Next.js)

- Next.js latest stable, **App Router**, TypeScript, **Tailwind CSS** (+ shadcn/ui for admin forms/dialogs; public site is fully custom per the mockups).
- **Public site**: dynamic routes `app/[[...slug]]/page.tsx` fetch the published page from the API and render via a **block registry** — `Record<BlockType, React.Component>` mapping 1:1 to the catalog in design-spec.md. On-demand revalidation (tag-based) triggered by the API's publish webhook; falls back to a short TTL.
- **Admin (`/area-reservada`)**: client-heavy React (forms, drag-and-drop via `dnd-kit`, block editor panels). Talks to the Go API via the `/api` rewrite; session cookie is first-party. A middleware guard redirects unauthenticated users to login.
- **Preview**: draft mode route that fetches the draft revision with the editor's session — same block registry renders it.

## 5. Development environment

- This machine has Postgres running locally on :5432 (a previous `airfa_dev` DB / `airfa` role may exist from the discarded attempt — **drop and recreate cleanly in Phase 0**).
- `docker-compose.dev.yml` optional; native Postgres is fine for dev.
- Makefile targets (Phase 0): `make dev` (api + web concurrently), `make migrate`, `make seed`, `make sqlc`, `make test`.
- Seed script (Phase 1) creates the admin user and the initial airfa.pt content per `content-inventory.md`.

## 6. Deployment (Phase 3)

- Dockerfiles: `api/` (distroless/scratch static Go binary) and `web/` (Next standalone output).
- Dokploy on the owner's VPS: 3 services (api, web, postgres) + volumes (`pgdata`, `media`); Dokploy handles TLS/domain routing. Web is the public entrypoint; it proxies `/api/*` and `/media/*` to the api service on the internal network.
- Backups: nightly `pg_dump` + media volume snapshot (Dokploy scheduled task or cron container), retained ≥14 days.

## 7. Testing strategy

- Go: unit tests for auth/versioning/publish logic; handler tests against a test DB (dockertest or a dedicated test database).
- Web: component tests for the block registry (each block renders its fixture), Playwright smoke: public pages render, login works, edit→publish→public update round-trip.
- CI later (nothing exotic — `go test`, `next build`, `tsc --noEmit`).
