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
- [x] three migration pairs created and filled
- [x] `make migrate` applies cleanly; `\dt` shows all 11 tables
- [x] `make migrate-down` works (down files valid), then re-applied

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

**`apps/api/queries/collections.sql`** — ✅ complete as of 2026-07-16 (all four collections; verified against the generated code).

> **Rule: a list query never `select`s a jsonb column.** The handlers serialize sqlc's row structs straight to JSON, and jsonb arrives in Go as `[]byte`, which `encoding/json` renders as **base64 garbage** rather than JSON. So `events` omits `body` and `activities` omits `info` — lists are flat summaries. Full bodies come from per-item detail endpoints later (Phase 2), which map jsonb through `json.RawMessage`.

```sql
-- name: ListPublishedEvents :many
select id, title, poster_media_id, starts_at, ends_at, sort
from events
where status = 'published'
order by starts_at desc nulls last, sort asc
limit $1 offset $2;

-- name: CountPublishedEvents :one
select count(*) from events where status = 'published';

-- name: ListPublishedPosts :many
select id, title, slug, cover_media_id, excerpt, published_at
from posts
where status = 'published'
order by published_at desc nulls last, created_at desc
limit $1 offset $2;

-- name: CountPublishedPosts :one
select count(*) from posts where status = 'published';

-- name: ListPartners :many
select id, name, logo_media_id, url, sort
from partners
order by sort asc, name asc
limit $1 offset $2;

-- name: CountPartners :one
select count(*) from partners;

-- name: ListPublishedActivities :many
select id, name, category, image_media_id, sort
from activities
where status = 'published'
order by sort asc, name asc
limit $1 offset $2;

-- name: CountPublishedActivities :one
select count(*) from activities where status = 'published';
```

> `partners` has no `status` column (see the A.1 schema) — every partner row is public, so there's no published filter and `CountPartners` counts all rows.

**Generate + verify:**
```bash
make sqlc
go build ./apps/api/...
```
- [x] `site.sql`, `pages.sql`, `collections.sql` written — all four collections
- [x] `make sqlc` generates `apps/api/internal/db/*` with no errors — 11 query funcs across `site/pages/collections.sql.go`
- [x] `go build ./apps/api/...` compiles

> **Generated types A.3 depends on** (verified 2026-07-16 — check here before writing handler code):
> `ListPublishedEventsRow`, `ListPublishedPostsRow`, `ListPublishedActivitiesRow`, and their `…Params` twins, plus `ListPartnersParams`.
> **There is no `ListPartnersRow`** — `ListPartners` selects every column of `partners`, so sqlc reuses the `db.Partner` model. Nullable `partners.url` is `*string` (`emit_pointers_for_null_types`); nullable UUID/timestamp columns stay `pgtype.UUID` / `pgtype.Timestamptz`. All marshal to correct JSON.

---

### A.3 — Public read endpoints
**Goal:** wire the DB pool and expose the three public surfaces (`agent/architecture.md §2.1`).

**Endpoints being built:**

| Method | Path | Returns |
|---|---|---|
| `GET` | `/v1/site` | `{ "settings": {…data…}, "menus": { "main":[…], "secondary":[…], "utility":[…], "footer":[…] } }` |
| `GET` | `/v1/pages/{slug}` | `{ "slug","title","seo":{…},"blocks":[…] }` — `404` envelope if not found/published |
| `GET` | `/v1/collections/{name}` | JSON array of published items; `name ∈ {events,posts,partners,activities}`; `?limit=&offset=`; sets `X-Total-Count` header |

**Design notes (read before pasting — these explain the non-obvious bits):**
- **`internal/http` uses `package httpapi`.** The directory name follows architecture.md §2, but a package literally named `http` sitting next to `net/http` imports is a readability trap. Go allows the mismatch; `main.go` imports it with an explicit alias.
- **jsonb → `json.RawMessage`, never `[]byte`.** sqlc returns jsonb columns (`seo`, `blocks`, `settings.data`, `menus.items`) as `[]byte`, and `encoding/json` base64-encodes `[]byte`. Every jsonb field is wrapped in `json.RawMessage` before it goes out, which splices the stored JSON in verbatim. Get this wrong and `/v1/pages/inicio` returns `"blocks":"W3siaWQi..."`.
- **Empty list ≠ null.** sqlc returns a nil slice when a query matches nothing, and `json.Marshal(nil slice)` emits `null`. The frontend must always get `[]`, so each branch normalizes nil → empty slice.
- **`/v1/site` always emits all four menu zones**, empty when the row is absent, so the frontend can index `menus.main` without a guard.
- **UUID/timestamp JSON is free.** `pgtype.UUID`, `pgtype.Timestamptz`, and `pgtype.Text` all implement `MarshalJSON` in pgx v5 — they emit `"uuid-string"` / RFC3339 / `null` correctly, so the sqlc row structs serialize directly. Their `emit_json_tags` snake_case names already match the API convention.

---

#### Step 1 — `apps/api/internal/db/pool.go` (new file)

Hand-written, lives beside the generated files. sqlc only overwrites files it generates, so it will not clobber this.

```go
package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool opens a pgx connection pool and verifies it is reachable.
func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}
```

#### Step 2 — `apps/api/internal/http/errors.go` (new file)

The error envelope from architecture.md §2.1. **Everything here is English** — messages, logs, comments. API errors are never rendered raw to users: the frontend switches on the stable `code` and shows its own pt-PT copy. Only text that literally appears on the page or in the CMS admin UI is Portuguese.

```go
package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
)

type apiError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type errorEnvelope struct {
	Error apiError `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write json response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorEnvelope{Error: apiError{Code: code, Message: message}})
}

// serverError logs the real cause and returns an opaque 500 to the client.
func serverError(w http.ResponseWriter, what string, err error) {
	log.Printf("error: %s: %v", what, err)
	writeError(w, http.StatusInternalServerError, "internal", "internal server error")
}
```

#### Step 3 — `apps/api/internal/http/router.go` (new file)

```go
package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/RMLoureiro/airfa-website-2.0/apps/api/internal/db"
)

// Server holds the dependencies shared by the handlers.
type Server struct {
	q *db.Queries
}

// Router builds the full HTTP surface: /healthz plus the public /v1 reads.
func Router(q *db.Queries) http.Handler {
	s := &Server{q: q}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/v1", func(r chi.Router) {
		r.Get("/site", s.getSite)
		r.Get("/pages/{slug}", s.getPage)
		r.Get("/collections/{name}", s.getCollection)
	})

	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusNotFound, "not_found", "resource not found")
	})

	return r
}
```

#### Step 4 — `apps/api/internal/http/site.go` (new file)

```go
package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
)

// menuZones are always present in the response, empty when unset, so the
// frontend can read menus.main without a nil check.
var menuZones = []string{"main", "secondary", "utility", "footer"}

// getSite handles GET /v1/site.
func (s *Server) getSite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	settings := json.RawMessage(`{}`)
	data, err := s.q.GetSettings(ctx)
	switch {
	case err == nil:
		settings = json.RawMessage(data)
	case errors.Is(err, pgx.ErrNoRows):
		// No settings row seeded yet — an empty object is a valid answer.
	default:
		serverError(w, "get settings", err)
		return
	}

	rows, err := s.q.ListMenus(ctx)
	if err != nil {
		serverError(w, "list menus", err)
		return
	}

	menus := make(map[string]json.RawMessage, len(menuZones))
	for _, zone := range menuZones {
		menus[zone] = json.RawMessage(`[]`)
	}
	for _, m := range rows {
		menus[m.Zone] = json.RawMessage(m.Items)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"settings": settings,
		"menus":    menus,
	})
}
```

#### Step 5 — `apps/api/internal/http/pages.go` (new file)

```go
package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

// getPage handles GET /v1/pages/{slug}. Only published pages are visible;
// an unpublished or missing slug is a 404 either way.
func (s *Server) getPage(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	page, err := s.q.GetPublishedPageBySlug(r.Context(), slug)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "page_not_found", "page not found")
		return
	}
	if err != nil {
		serverError(w, "get published page", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"slug":   page.Slug,
		"title":  page.Title,
		"seo":    json.RawMessage(page.Seo),
		"blocks": json.RawMessage(page.Blocks),
	})
}
```

#### Step 6 — `apps/api/internal/http/collections.go` (new file)

```go
package httpapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/RMLoureiro/airfa-website-2.0/apps/api/internal/db"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

// getCollection handles GET /v1/collections/{name}?limit=&offset=.
func (s *Server) getCollection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, offset, err := pagination(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_pagination", err.Error())
		return
	}

	var (
		items any
		total int64
	)

	switch name := chi.URLParam(r, "name"); name {
	case "events":
		rows, err := s.q.ListPublishedEvents(ctx, db.ListPublishedEventsParams{Limit: limit, Offset: offset})
		if err != nil {
			serverError(w, "list events", err)
			return
		}
		if rows == nil {
			rows = []db.ListPublishedEventsRow{}
		}
		if total, err = s.q.CountPublishedEvents(ctx); err != nil {
			serverError(w, "count events", err)
			return
		}
		items = rows

	case "posts":
		rows, err := s.q.ListPublishedPosts(ctx, db.ListPublishedPostsParams{Limit: limit, Offset: offset})
		if err != nil {
			serverError(w, "list posts", err)
			return
		}
		if rows == nil {
			rows = []db.ListPublishedPostsRow{}
		}
		if total, err = s.q.CountPublishedPosts(ctx); err != nil {
			serverError(w, "count posts", err)
			return
		}
		items = rows

	case "partners":
		rows, err := s.q.ListPartners(ctx, db.ListPartnersParams{Limit: limit, Offset: offset})
		if err != nil {
			serverError(w, "list partners", err)
			return
		}
		// Note: []db.Partner, not []db.ListPartnersRow — ListPartners selects every
		// column of the table, so sqlc reuses the Partner model instead of
		// generating a row type. There is no ListPartnersRow.
		if rows == nil {
			rows = []db.Partner{}
		}
		if total, err = s.q.CountPartners(ctx); err != nil {
			serverError(w, "count partners", err)
			return
		}
		items = rows

	case "activities":
		rows, err := s.q.ListPublishedActivities(ctx, db.ListPublishedActivitiesParams{Limit: limit, Offset: offset})
		if err != nil {
			serverError(w, "list activities", err)
			return
		}
		if rows == nil {
			rows = []db.ListPublishedActivitiesRow{}
		}
		if total, err = s.q.CountPublishedActivities(ctx); err != nil {
			serverError(w, "count activities", err)
			return
		}
		items = rows

	default:
		writeError(w, http.StatusNotFound, "unknown_collection",
			fmt.Sprintf("unknown collection: %q", name))
		return
	}

	w.Header().Set("X-Total-Count", strconv.FormatInt(total, 10))
	writeJSON(w, http.StatusOK, items)
}

// pagination reads ?limit=&offset= with defaults and an upper bound.
func pagination(r *http.Request) (limit, offset int32, err error) {
	limit, offset = defaultLimit, 0

	if v := r.URL.Query().Get("limit"); v != "" {
		n, convErr := strconv.Atoi(v)
		if convErr != nil || n < 1 {
			return 0, 0, fmt.Errorf("invalid limit parameter: %q", v)
		}
		if n > maxLimit {
			n = maxLimit
		}
		limit = int32(n)
	}

	if v := r.URL.Query().Get("offset"); v != "" {
		n, convErr := strconv.Atoi(v)
		if convErr != nil || n < 0 {
			return 0, 0, fmt.Errorf("invalid offset parameter: %q", v)
		}
		offset = int32(n)
	}

	return limit, offset, nil
}
```

> The four branches repeat because each sqlc row type is distinct and Go generics can't tidy this without more machinery than it saves. Leave it explicit.

#### Step 7 — `apps/api/cmd/server/main.go` (replace the whole file)

Adds the pool, the router, and graceful shutdown. Note it now **fails fast if `DATABASE_URL` is unset** — see Step 8.

```go
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RMLoureiro/airfa-website-2.0/apps/api/internal/db"
	httpapi "github.com/RMLoureiro/airfa-website-2.0/apps/api/internal/http"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return errors.New("DATABASE_URL is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.NewPool(ctx, dsn)
	if err != nil {
		return err
	}
	defer pool.Close()

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           httpapi.Router(db.New(pool)),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	log.Printf("api listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
```

#### Step 8 — Make `DATABASE_URL` reach `make dev` (edit `Makefile`)

**Why:** the `dev` target runs `go run ./cmd/server` in a subshell. The Makefile defines `DATABASE_URL` but never exports it, so from Step 7 onward **`make dev` would die with "DATABASE_URL is not set"**. One line fixes it — no `.env` loader, no hardcoded credentials in Go.

Add `export DATABASE_URL` directly under the existing assignment at the top of the root `Makefile`:

```makefile
DATABASE_URL ?= postgres://airfa:airfa@localhost:5432/airfa_dev?sslmode=disable
export DATABASE_URL
```

> ⚠️ **The DSN must stay on ONE line.** This broke once (fixed 2026-07-17): the line was split at `airfa_dev?` / `sslmode=disable`. Make does not error — it reads `sslmode=disable` as a *separate variable* and leaves `DATABASE_URL` ending in a bare `?`, which then fails at connect time with a confusing SSL error. Check it with:
> ```bash
> make -pn 2>/dev/null | grep -E '^(DATABASE_URL|sslmode)'
> # must print exactly one line, ending in ?sslmode=disable — and no `sslmode = disable` line
> ```

#### Step 9 — Resolve dependencies and build

`pgxpool` pulls in two modules (`puddle`, `golang.org/x/sync`) that aren't in `go.mod` yet, because nothing had imported the pool until now.

```bash
cd apps/api
go mod tidy
go build ./...
go vet ./...
```

- [x] Step 1 — `internal/db/pool.go` created
- [x] Steps 2–6 — `internal/http/{errors,router,site,pages,collections}.go` created
- [x] Step 7 — `cmd/server/main.go` replaced
- [x] Step 8 — `export DATABASE_URL` added to the Makefile *(the DSN got split across two lines; found + fixed 2026-07-17 — see the warning above)*
- [x] Step 9 — `go mod tidy` && `go build ./...` && `go vet ./...` all clean — **verified 2026-07-17, both exit 0**

---

#### Step 10 — Seed a little data and verify

```bash
psql "postgres://airfa:airfa@localhost:5432/airfa_dev?sslmode=disable" -c \
"insert into pages(slug,title,status) values('inicio','Início','published');
 insert into page_revisions(page_id,kind,blocks)
 select id,'published','[{\"id\":\"b1\",\"type\":\"rich-text\",\"data\":{\"html\":\"<p>olá</p>\"}}]'::jsonb from pages where slug='inicio';
 insert into menus(zone,items) values('main','[]'),('secondary','[]'),('utility','[]'),('footer','[]');
 insert into settings(data) values('{}');"
```

Then, in one terminal:

```bash
make dev     # or: cd apps/api && go run ./cmd/server
```

And in another:

```bash
curl -s localhost:8080/healthz | jq
curl -s localhost:8080/v1/site | jq
curl -s localhost:8080/v1/pages/inicio | jq
curl -s "localhost:8080/v1/collections/events?limit=10" -i     # check the X-Total-Count header
curl -s localhost:8080/v1/collections/posts | jq               # expect [] — NOT null
curl -s localhost:8080/v1/collections/partners | jq
curl -s localhost:8080/v1/collections/activities | jq
curl -s localhost:8080/v1/collections/nope | jq                # → 404 unknown_collection
curl -s "localhost:8080/v1/collections/events?limit=abc" | jq  # → 400 invalid_pagination
curl -s -o /dev/null -w "%{http_code}\n" localhost:8080/v1/pages/does-not-exist   # → 404
```

**What "correct" looks like** — on `/v1/pages/inicio`, `blocks` must be a real JSON **array** and `seo` a real JSON **object**:

```json
{
  "slug": "inicio",
  "title": "Início",
  "seo": {},
  "blocks": [{ "id": "b1", "type": "rich-text", "data": { "html": "<p>olá</p>" } }]
}
```

If `blocks` comes back as a base64 string (`"W3siaWQi..."`), a `json.RawMessage` wrap was missed — see the design notes above.

- [x] DB pool wired from `DATABASE_URL` in `main.go`
- [x] `GET /v1/site` returns settings + all four menu zones
- [x] `GET /v1/pages/inicio` returns the block tree as real JSON (not base64)
- [x] `GET /v1/collections/events` returns an array + `X-Total-Count`
- [x] the other three collections return `[]` (not `null`) when empty
- [x] unknown slug → `404`; unknown collection → `404`; bad `limit` → `400` — all with the error envelope

---

## ✅ Definition of done for this phase

- [x] `make dev` starts API + web together
- [x] All 11 tables migrate up **and** down cleanly
- [x] `sqlc` generates; `go build ./apps/api/...` compiles
- [x] `/v1/site`, `/v1/pages/:slug`, `/v1/collections/*` all return valid JSON
- [x] (optional) local checkpoint commit — plain message, no co-author line

**Next phase after this:** Phase 1 · Stage B — design tokens (Inter via `next/font`) + Tailwind colors sampled from the mockups.

---

## Decisions recap (context for this phase)

- **Monorepo:** `apps/` + `packages/`; Make is the unified runner; `go.work` for Go; npm workspaces for JS. (architecture.md §1.1)
- **Fonts:** Inter everywhere, self-hosted via `next/font` (Stage B).
- **Blog & Events:** both are real content — build/seed fully.
- **DOAR button:** links to an on-site IBAN / how-to-donate page.
- **Main nav label:** "Actividades".
