# Current Phase — Foundation + Phase 1 · Stage A (Data foundation)

> Working checklist for the active phase. Tick boxes as you complete them. Each step has **Goal → Do → Verify**. Run the Verify before moving on. When a whole section is done, run its "Section done when" check.
>
> Paths assume repo root = `/home/loureiro/dev/airfa-website-2.0`. The monorepo layout is `apps/api` (Go), `apps/web` (Next.js), `packages/shared` (shared TS). Dev DB DSN: `postgres://airfa:airfa@localhost:5432/airfa_dev?sslmode=disable`.

---

## ✅ Foundation complete — verified 2026-07-13

All of Step 0 (the roadmap's "Phase 0" scaffolding) is done and green. **The next agent starts at Phase 1 · Stage A** (further down this file).

**Scaffolding & tooling — done:**
- [x] `airfa_dev` dropped/recreated clean (discarded Prisma tables removed)
- [x] Go API at `apps/api` (module `github.com/RMLoureiro/airfa-website-2.0/apps/api`) — chi + `middleware.Logger` + `GET /healthz` → `{"status":"ok"}`
- [x] `go.work` at root (`use ./apps/api`)
- [x] Next.js app at `apps/web` (Next 16, TS, Tailwind v4); `/api/:path*` → Go API rewrite in `next.config.ts`
- [x] `packages/shared` stub; root npm workspace; `@airfa/shared` linked into `apps/web`
- [x] Root `Makefile` (dev/migrate/sqlc/seed/test), `.editorconfig`, `.gitignore`, `.env.example` in both apps
- [x] `migrate` + `sqlc v1.31.1` installed on PATH; `apps/api/sqlc.yaml` present
- [x] Checks green: `go vet`, `go build`, `go test`, `tsc --noEmit`, `next build`

**Deliberately NOT done — rolled into Stage A below (owner's call, 2026-07-13):**
- First migration (users, sessions) → **Stage A.1**
- sqlc actually generating → **Stage A.2**
- Web page fetching data from the API (the Phase 0 round-trip "demo") → **Stage A.3**

> These three are roadmap "Phase 0" items; the owner chose to fold them into Stage A rather than ship a separate mini-milestone. So Phase 0's *scaffolding* is complete; its *round-trip demo* lands naturally when Stage A's endpoints render.

**Environment facts (this machine — you'll need these):**
- Dev DB DSN: `postgres://airfa:airfa@localhost:5432/airfa_dev?sslmode=disable` (role `airfa`, password `airfa`). Postgres 18.
- `~/go/bin` was added to PATH via `~/.zshrc`, so `migrate`/`sqlc` resolve in **new** shells (Go 1.26.4, Node 26.4).
- Root `package.json` has `overrides.postcss: ^8.5.10` — **keep it**; it patches a postcss advisory bundled by Next. `npm ls` prints a benign `invalid: "8.4.31" from node_modules/next` note; that's expected.
- Nothing is committed yet — `apps/`, `packages/`, `go.work`, root `package.json`, `Makefile`, etc. are all untracked. First checkpoint commit is the owner's call.

---

## ▶ Foundation — remaining

### 0.1b — Confirm the dev DB password is set
**Goal:** the `airfa` role has a known password (needed from Stage A onward).
**Do (once, if you haven't):**
```bash
psql -U postgres -h localhost -p 5432 -c "ALTER ROLE airfa WITH LOGIN PASSWORD 'airfa';"
```
**Verify:**
```bash
psql "postgres://airfa:airfa@localhost:5432/airfa_dev?sslmode=disable" -c "select current_user;"
```
- [x] `airfa` login works with the DSN above

---

### 0.4 — Next.js web app at `apps/web`
**Goal:** Next.js + TS + Tailwind app that proxies `/api/*` to the Go API.
**Do (from repo root):**
```bash
npx create-next-app@latest apps/web --ts --tailwind --app --src-dir --eslint --use-npm --no-turbopack
```
Then edit `apps/web/next.config.ts` — add the `rewrites` inside the config object:
```ts
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: `${process.env.API_URL ?? "http://localhost:8080"}/v1/:path*`,
      },
    ];
  },
};

export default nextConfig;
```
**Verify:**
```bash
cd apps/web && npm run dev     # http://localhost:3000 renders the Next.js starter
```
- [x] `apps/web` created
- [x] `next.config.ts` has the `/api/:path*` rewrite
- [x] `localhost:3000` renders

---

### 0.5 — Shared package `packages/shared`
**Goal:** home for shared TS (block types + JSON fixtures, used in Stage D). Stub for now.
**Do (from repo root):**
```bash
mkdir -p packages/shared/src
```
Create `packages/shared/package.json`:
```json
{
  "name": "@airfa/shared",
  "version": "0.0.0",
  "private": true,
  "type": "module",
  "main": "src/index.ts",
  "exports": { ".": "./src/index.ts" }
}
```
Create `packages/shared/src/index.ts`:
```ts
// Shared TypeScript for the AIRFA monorepo.
// Block types + JSON fixtures will live here (Phase 1 · Stage D).
export const SHARED = true;
```
- [x] `packages/shared/package.json` created
- [x] `packages/shared/src/index.ts` created

---

### 0.6 — Root npm workspace
**Goal:** one root workspace so `apps/web` can import `@airfa/shared`, with hoisted deps.
**Do (from repo root):**
Create root `package.json`:
```json
{
  "name": "airfa-website-2.0",
  "private": true,
  "workspaces": ["apps/web", "packages/*"]
}
```
Then hoist and link:
```bash
rm -f apps/web/package-lock.json
rm -rf apps/web/node_modules
npm install
npm install @airfa/shared --workspace apps/web
```
**Verify:**
```bash
ls -l node_modules/@airfa/shared     # should be a symlink into packages/shared
cd apps/web && npm run dev           # still works
```
- [x] root `package.json` with `workspaces`
- [x] `npm install` hoisted deps to root `node_modules`
- [x] `@airfa/shared` linked as a dependency of `apps/web` (symlink present)

---

### 0.7 — Root tooling (Makefile, editorconfig, gitignore, env examples)
**Goal:** one command runs both apps; consistent editor/ignore config; documented env.

**`Makefile`** at root (⚠️ recipe lines must be **tabs**, not spaces):
```makefile
DATABASE_URL ?= postgres://airfa:airfa@localhost:5432/airfa_dev?sslmode=disable

.PHONY: dev migrate migrate-down migrate-create sqlc seed test

dev:
	@trap 'kill 0' EXIT; \
	(cd apps/api && go run ./cmd/server) & \
	(npm run dev --workspace apps/web) & \
	wait

migrate:
	migrate -path apps/api/migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path apps/api/migrations -database "$(DATABASE_URL)" down 1

# usage: make migrate-create name=add_something
migrate-create:
	migrate create -ext sql -dir apps/api/migrations -seq $(name)

sqlc:
	cd apps/api && sqlc generate

seed:
	cd apps/api && go run ./cmd/seed

test:
	cd apps/api && go test ./...
	npm test --workspaces --if-present
```

**`.editorconfig`** at root:
```ini
root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true
indent_style = space
indent_size = 2

[*.go]
indent_style = tab

[Makefile]
indent_style = tab
```

**`.gitignore`** — append these lines to the existing file:
```gitignore
# build artifacts
/apps/api/server
/apps/api/bin/
# next
/apps/web/.next/
# local env
*.env.local
.env.local
```

**`apps/api/.env.example`**:
```dotenv
DATABASE_URL=postgres://airfa:airfa@localhost:5432/airfa_dev?sslmode=disable
PORT=8080
MEDIA_DIR=./.data/media
SESSION_TTL=720h
```

**`apps/web/.env.example`**:
```dotenv
API_URL=http://localhost:8080
```

**Verify:**
```bash
make dev     # brings up BOTH apps; Ctrl-C stops both
# in another terminal:
curl -s localhost:8080/healthz    # {"status":"ok"}
# browser: localhost:3000 renders
```
- [x] `Makefile` (tabs verified — `make dev` works)
- [x] `.editorconfig`
- [x] `.gitignore` updated
- [x] `apps/api/.env.example` and `apps/web/.env.example`
- [x] `make dev` runs both apps together

---

### 0.8 — Migration + sqlc tooling
**Goal:** install the schema tool and the query codegen; configure sqlc.
**Do:**
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```
> These install to `$(go env GOPATH)/bin` (usually `~/go/bin`). Ensure that's on your `PATH` (`echo $PATH | tr ':' '\n' | grep go/bin`). If not, add `export PATH="$PATH:$(go env GOPATH)/bin"` to your shell profile.

Create `apps/api/sqlc.yaml`:
```yaml
version: "2"
sql:
  - engine: "postgresql"
    schema: "migrations"
    queries: "queries"
    gen:
      go:
        package: "db"
        out: "internal/db"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_pointers_for_null_types: true
```
**Verify:**
```bash
migrate -version     # prints a version
sqlc version         # prints a version
```
- [x] `migrate` installed and on PATH
- [x] `sqlc` installed and on PATH
- [x] `apps/api/sqlc.yaml` created

---

### 0.9 — Cleanup
**Goal:** remove the stray compiled binary left from an earlier `go build`.
**Do:**
```bash
rm -f apps/api/server
```
- [x] `apps/api/server` binary removed (it's now git-ignored via 0.7)

---

> **Foundation done when:** `make dev` runs both apps, `apps/web` imports `@airfa/shared` cleanly, `curl localhost:8080/healthz` works, and `migrate`/`sqlc` print versions.

---

## ▶ Phase 1 · Stage A — Data foundation (schema + read API)

Pure backend. No visuals yet — you're building the tables and the read endpoints the site will consume. All API JSON is `snake_case`, no auth.

### A.1 — Migrations (schema)
**Goal:** create the tables from `agent/architecture.md §3` in three logical migrations.

**Create the migration files:**
```bash
make migrate-create name=auth
make migrate-create name=content
make migrate-create name=collections
```
This makes `000001_auth.up.sql`/`.down.sql`, `000002_content.*`, `000003_collections.*` under `apps/api/migrations/`. Paste the SQL below into each. (Postgres 18 has `gen_random_uuid()` built in — no extension needed.)

**`000001_auth.up.sql`**
```sql
create table users (
  id            uuid primary key default gen_random_uuid(),
  name          text not null,
  email         text not null unique,
  password_hash text not null,
  role          text not null default 'editor' check (role in ('admin','editor')),
  created_at    timestamptz not null default now(),
  disabled_at   timestamptz
);

create table sessions (
  id         uuid primary key default gen_random_uuid(),
  user_id    uuid not null references users(id) on delete cascade,
  token_hash text not null unique,
  expires_at timestamptz not null,
  created_at timestamptz not null default now()
);
create index sessions_user_id_idx on sessions (user_id);
```
**`000001_auth.down.sql`**
```sql
drop table if exists sessions;
drop table if exists users;
```

**`000002_content.up.sql`**
```sql
create table pages (
  id         uuid primary key default gen_random_uuid(),
  slug       text not null unique,
  title      text not null,
  seo        jsonb not null default '{}'::jsonb,
  parent_id  uuid references pages(id) on delete set null,
  sort       int not null default 0,
  status     text not null default 'draft' check (status in ('draft','published')),
  deleted_at timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table page_revisions (
  id         uuid primary key default gen_random_uuid(),
  page_id    uuid not null references pages(id) on delete cascade,
  kind       text not null check (kind in ('draft','published','snapshot')),
  blocks     jsonb not null default '[]'::jsonb,
  seo        jsonb not null default '{}'::jsonb,
  created_by uuid references users(id) on delete set null,
  created_at timestamptz not null default now()
);
create index page_revisions_page_kind_idx on page_revisions (page_id, kind);

create table menus (
  id         uuid primary key default gen_random_uuid(),
  zone       text not null unique check (zone in ('main','secondary','utility','footer')),
  items      jsonb not null default '[]'::jsonb,
  updated_at timestamptz not null default now()
);

create table settings (
  id   uuid primary key default gen_random_uuid(),
  data jsonb not null default '{}'::jsonb
);

create table media (
  id         uuid primary key default gen_random_uuid(),
  kind       text not null check (kind in ('image','pdf')),
  path       text not null,
  filename   text not null,
  size       bigint not null default 0,
  width      int,
  height     int,
  alt        text not null default '',
  created_at timestamptz not null default now()
);
```
**`000002_content.down.sql`**
```sql
drop table if exists media;
drop table if exists settings;
drop table if exists menus;
drop table if exists page_revisions;
drop table if exists pages;
```

**`000003_collections.up.sql`**
```sql
create table events (
  id              uuid primary key default gen_random_uuid(),
  title           text not null,
  poster_media_id uuid references media(id) on delete set null,
  starts_at       timestamptz,
  ends_at         timestamptz,
  body            jsonb not null default '{}'::jsonb,
  status          text not null default 'draft' check (status in ('draft','published')),
  sort            int not null default 0,
  created_at      timestamptz not null default now(),
  updated_at      timestamptz not null default now()
);

create table posts (
  id             uuid primary key default gen_random_uuid(),
  title          text not null,
  slug           text not null unique,
  cover_media_id uuid references media(id) on delete set null,
  excerpt        text not null default '',
  blocks         jsonb not null default '[]'::jsonb,
  published_at   timestamptz,
  status         text not null default 'draft' check (status in ('draft','published')),
  created_at     timestamptz not null default now(),
  updated_at     timestamptz not null default now()
);

create table partners (
  id            uuid primary key default gen_random_uuid(),
  name          text not null,
  logo_media_id uuid references media(id) on delete set null,
  url           text,
  sort          int not null default 0
);

create table activities (
  id             uuid primary key default gen_random_uuid(),
  name           text not null,
  category       text not null default '',
  image_media_id uuid references media(id) on delete set null,
  info           jsonb not null default '{}'::jsonb,
  sort           int not null default 0,
  status         text not null default 'draft' check (status in ('draft','published')),
  created_at     timestamptz not null default now(),
  updated_at     timestamptz not null default now()
);
```
**`000003_collections.down.sql`**
```sql
drop table if exists activities;
drop table if exists partners;
drop table if exists posts;
drop table if exists events;
```

**Apply + verify:**
```bash
make migrate
psql "postgres://airfa:airfa@localhost:5432/airfa_dev?sslmode=disable" -c "\dt"
# then test rollback works, and re-apply:
make migrate-down    # (drops last migration; run 3x to fully unwind if you want)
make migrate
```
- [ ] three migration pairs created and filled
- [ ] `make migrate` applies cleanly; `\dt` shows all 11 tables
- [ ] `make migrate-down` works (down files valid), then re-applied

---

### A.2 — Queries + sqlc generate
**Goal:** type-safe Go for the reads Stage A needs. Put `.sql` files in `apps/api/queries/`.

**`apps/api/queries/site.sql`**
```sql
-- name: GetSettings :one
select data from settings limit 1;

-- name: ListMenus :many
select zone, items from menus;
```

**`apps/api/queries/pages.sql`**
```sql
-- name: GetPublishedPageBySlug :one
select p.slug, p.title, p.seo, r.blocks
from pages p
join page_revisions r on r.page_id = p.id and r.kind = 'published'
where p.slug = $1 and p.deleted_at is null;
```

**`apps/api/queries/collections.sql`** (events shown; posts/partners/activities follow the same shape)
```sql
-- name: ListPublishedEvents :many
select id, title, poster_media_id, starts_at, ends_at, sort
from events
where status = 'published'
order by starts_at desc nulls last, sort asc
limit $1 offset $2;

-- name: CountPublishedEvents :one
select count(*) from events where status = 'published';
```

**Generate + verify:**
```bash
make sqlc
go build ./apps/api/...
```
- [ ] query files written (`site.sql`, `pages.sql`, `collections.sql` with all four collections)
- [ ] `make sqlc` generates `apps/api/internal/db/*` with no errors
- [ ] `go build ./apps/api/...` compiles

---

### A.3 — Public read endpoints
**Goal:** wire the DB pool and expose the three public surfaces (`agent/architecture.md §2.1`).

**Wiring (spec):**
- `apps/api/internal/db` — a `func NewPool(ctx, dsn) (*pgxpool.Pool, error)` opening the pool from `DATABASE_URL`; sqlc's generated `db.New(pool)` gives a `*db.Queries`.
- `apps/api/internal/http` — handlers holding a `*db.Queries`; a `Router(q *db.Queries) http.Handler` that mounts `/healthz` + the `/v1/*` routes below.
- `apps/api/cmd/server/main.go` — read env (`DATABASE_URL`, `PORT`), open the pool, build queries, mount the router, listen.
- Errors use the envelope from architecture.md §2.1: `{"error":{"code","message","fields"}}`.

**Endpoints:**

| Method | Path | Returns |
|---|---|---|
| `GET` | `/v1/site` | `{ "settings": {…data…}, "menus": { "main":[…], "secondary":[…], "utility":[…], "footer":[…] } }` |
| `GET` | `/v1/pages/{slug}` | `{ "slug","title","seo":{…},"blocks":[…] }` — `404` envelope if not found/published |
| `GET` | `/v1/collections/{name}` | JSON array of published items; `name ∈ {events,posts,partners,activities}`; `?limit=&offset=`; sets `X-Total-Count` header |

**Verify (insert one row by hand first so there's data):**
```bash
psql "postgres://airfa:airfa@localhost:5432/airfa_dev?sslmode=disable" -c \
"insert into pages(slug,title,status) values('inicio','Início','published');
 insert into page_revisions(page_id,kind,blocks)
 select id,'published','[{\"id\":\"b1\",\"type\":\"rich-text\",\"data\":{\"html\":\"<p>olá</p>\"}}]'::jsonb from pages where slug='inicio';
 insert into menus(zone,items) values('main','[]'),('secondary','[]'),('utility','[]'),('footer','[]');
 insert into settings(data) values('{}');"

make dev   # (or run the api alone)
curl -s localhost:8080/v1/site | jq
curl -s localhost:8080/v1/pages/inicio | jq
curl -s "localhost:8080/v1/collections/events?limit=10" -i   # check X-Total-Count header
curl -s -o /dev/null -w "%{http_code}\n" localhost:8080/v1/pages/does-not-exist   # → 404
```
- [ ] DB pool wired from `DATABASE_URL` in `main.go`
- [ ] `GET /v1/site` returns settings + menus
- [ ] `GET /v1/pages/inicio` returns the block tree
- [ ] `GET /v1/collections/events` returns array + `X-Total-Count`
- [ ] unknown slug → `404` with error envelope

> `main.go`, the db pool, and the handlers are code (not just config). If you want the full source for any of them rather than the spec, ask and it'll be provided file-by-file.

---

## ✅ Definition of done for this phase

- [ ] `make dev` starts API + web together
- [ ] All 11 tables migrate up **and** down cleanly
- [ ] `sqlc` generates; `go build ./apps/api/...` compiles
- [ ] `/v1/site`, `/v1/pages/:slug`, `/v1/collections/*` all return valid JSON
- [ ] (optional) local checkpoint commit — plain message, no co-author line

**Next phase after this:** Phase 1 · Stage B — design tokens (Inter via `next/font`) + Tailwind colors sampled from the mockups.

---

## Decisions recap (context for this phase)

- **Monorepo:** `apps/` + `packages/`; Make is the unified runner; `go.work` for Go; npm workspaces for JS. (architecture.md §1.1)
- **Fonts:** Inter everywhere, self-hosted via `next/font` (Stage B).
- **Blog & Events:** both are real content — build/seed fully.
- **DOAR button:** links to an on-site IBAN / how-to-donate page.
- **Main nav label:** "Actividades".
