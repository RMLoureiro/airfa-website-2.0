# Architecture

> **Living document.** If implementation deviates from anything here, update this file in the same session and log the reason in `STATUS.md → Decisions`. Sections marked *(Phase 0/1/2/3)* are commitments, not history — they describe what to build.

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

## 1.1 Repository layout (monorepo)

Single Git repo, polyglot (Go + TypeScript), organized with the `apps/` + `packages/` convention. **The root `Makefile` is the one task runner across both languages** (npm workspaces / Turborepo are JS-only, so Make stays the top-level entrypoint).

```
airfa-website-2.0/
  apps/
    api/                 Go service (module github.com/RMLoureiro/airfa-website-2.0/apps/api)
    web/                 Next.js app (npm workspace)
  packages/
    shared/              shared TS: API DTO types, block TS types + JSON fixtures (block-tree sync, §3.1)
  agent/                 project docs (specs, design, roadmap)
  go.work                Go workspace (points at apps/api; ready for more Go modules)
  package.json           root — npm workspaces: ["apps/web", "packages/*"]
  Makefile               unified runner: dev / migrate / sqlc / seed / test
  .editorconfig .gitignore
```

- **JS side:** npm workspaces at the root; `apps/web` and everything under `packages/*` are workspaces, so shared TS (e.g. block types + fixtures) is imported by path, versioned once. Turborepo can be layered on later for JS build caching — not needed while there's a single JS build target.
- **Go side:** one module at `apps/api` under a root `go.work`. The workspace file makes adding a second Go module (e.g. a CLI or worker) a one-line change without import-path churn.
- **Shared block catalog** (§3.1) lives in `packages/shared` so the TS block types and JSON fixtures have a single home the web app consumes; the Go structs mirror them and the fixtures test both sides for drift.

## 2. Backend (Go)

- **Go** latest stable; standard project layout (under `apps/api/`):
  ```
  apps/api/
    cmd/server/main.go
    internal/
      http/        handlers, middleware, router (chi)
      auth/        sessions, password hashing (argon2id), RBAC
      content/     pages, blocks, menus, collections, versions
      media/       upload, image variants, static serving
      db/          sqlc-generated code, pgx pool
    migrations/    golang-migrate SQL files
    queries/       sqlc source queries
    sqlc.yaml
    go.mod
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

### 2.1 API conventions

- Base path `/v1`; JSON everywhere; `snake_case` field names in the API (frontend maps to camelCase at its boundary).
- Errors: consistent envelope `{ "error": { "code": "string", "message": "human readable (English)", "fields": { "field": "problem" } } }` with proper status codes (400 validation, 401 unauthenticated, 403 forbidden, 404, 409 conflict e.g. duplicate slug, 422 semantic).
- **API error `message` is English, always** (decided 2026-07-17). API errors are never shown to users raw — the frontend keys off the stable `code` and renders its own pt-PT copy. `message` and `fields` are for developers and logs. Anything user-facing gets translated at the frontend boundary, which is also where the CMS admin UI's Portuguese strings live. Corollary: **`code` is the contract** — it must be stable and machine-readable; never make the frontend match on `message` text.
- Lists: `?limit=&offset=` + `X-Total-Count` header (collections are small; no cursor pagination needed).
- Auth: session cookie (`airfa_session`, httpOnly, Secure, SameSite=Lax). CSRF: state-changing admin routes require `X-Requested-With` or a double-submit token — decide in Phase 2 and record here.
- All admin mutations are logged (who, what, when) — a simple `audit_log` table; cheap and invaluable for a multi-editor CMS.
- Media uploads: `multipart/form-data`, size limits (images ≤ 10 MB, PDFs ≤ 25 MB), MIME sniffed server-side, never trust the extension.

### 2.2 Environment variables

| Var | Service | Purpose |
|---|---|---|
| `DATABASE_URL` | api | Postgres DSN |
| `MEDIA_DIR` | api | Filesystem root for uploads (default `/data/media`) |
| `SESSION_TTL` | api | Session lifetime (default 720h) |
| `REVALIDATE_URL` / `REVALIDATE_SECRET` | api | Next.js on-publish webhook |
| `PORT` | both | Listen port |
| `API_URL` | web | Internal URL of the Go API (server-side fetches + `/api` rewrite target) |

Keep `.env.example` files in `api/` and `web/` current whenever a variable is added.

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
- **Homepage routing:** the page with reserved slug `inicio` renders at `/` (redirect `/inicio` → `/`). The slug is reserved: it cannot be deleted or renamed from the CMS, only edited.
- **Image variants** generated on upload: `thumb` 320px, `medium` 960px, `large` 1920px (longest edge, WebP + original format). The API returns all variant URLs with the media object; the frontend picks per context.
- A future `locale` column on `page_revisions`/collections is the multi-language path; don't build it now, don't block it.

### 3.1 Block tree shape (canonical example)

The `blocks` JSONB on a revision is an ordered array; every block has `id`, `type`, `data`. Nested item lists live inside `data`. Example (trimmed homepage):

```json
[
  {
    "id": "b1f4…",
    "type": "hero-slider",
    "data": {
      "slides": [
        {
          "id": "s1…",
          "image": { "mediaId": "m42…", "alt": "Palco do Cine-Teatro" },
          "headline": "Um espetáculo mais do que inesquecível",
          "underlinedWord": "inesquecível",
          "ctas": [
            { "label": "SABER MAIS", "variant": "outline", "link": { "type": "page", "pageId": "p7…" } },
            { "label": "SABER MAIS", "variant": "solid",   "link": { "type": "url", "url": "https://…" } }
          ]
        }
      ],
      "newsItems": [
        { "id": "n1…", "title": "O Xeque-Mate da Cultura Pop", "excerpt": "Como uma jogada antiga…", "link": { "type": "page", "pageId": "p9…" } }
      ]
    }
  },
  { "id": "b2…", "type": "section-band", "data": { "text": "PRÓXIMOS EVENTOS EM CARTAZ", "color": "red", "align": "center" } },
  { "id": "b3…", "type": "events-fan-carousel", "data": { "source": "auto", "limit": 8 } }
]
```

Rules:
- `id`s are UUIDs generated at creation and **never change** (drag-reorder moves objects, doesn't recreate them) — this keeps version diffs and the builder stable.
- Media is always referenced as `{ "mediaId", "alt" }` — never raw URLs — so replacing a file in the library updates every usage.
- Links are always the `{ "type": "page" | "url", … }` union described above.
- Each block type has a Go struct (validation) and a TS type + React component (rendering). The three MUST stay in sync — a shared JSON fixture per block type under `packages/shared` (block types + fixtures) is rendered in `apps/web` component tests to catch drift.

## 4. Frontend (Next.js)

- Next.js latest stable, **App Router**, TypeScript, **Tailwind CSS** (+ shadcn/ui for admin forms/dialogs; public site is fully custom per the mockups). Lives at `apps/web`.
- **Public site**: dynamic routes `apps/web/src/app/[[...slug]]/page.tsx` fetch the published page from the API and render via a **block registry** — `Record<BlockType, React.Component>` mapping 1:1 to the catalog in design-spec.md. On-demand revalidation (tag-based) triggered by the API's publish webhook; falls back to a short TTL.
- **Admin (`/area-reservada`)**: client-heavy React (forms, drag-and-drop via `dnd-kit`, block editor panels). Talks to the Go API via the `/api` rewrite; session cookie is first-party. A middleware guard redirects unauthenticated users to login.
- **Preview**: draft mode route that fetches the draft revision with the editor's session — same block registry renders it.

## 5. Development environment

- Toolchain verified on this machine (2026-07-12): Go 1.26.x, Node 26.x — Phase 0 needs no installs beyond project deps.
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
