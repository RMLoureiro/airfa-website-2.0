# Roadmap

Tick items as they complete; keep in sync with STATUS.md. Phases ship in order; each ends with something the owner can see/verify.

## Phase 0 — Foundations *(not started)*

Goal: empty but running skeleton; `make dev` brings everything up.

- [ ] Repo scaffolding: `api/` (Go module, chi server, healthcheck route), `web/` (create-next-app, TS, Tailwind), root `Makefile`, `.editorconfig`, `.gitignore`
- [ ] Postgres dev database created fresh (`airfa_dev`); `.env.example` for both apps
- [ ] Migrations tooling wired (golang-migrate) + first migration (users, sessions)
- [ ] sqlc configured and generating
- [ ] `web` → `api` proxy rewrite working (`/api/*`); web renders a page with data fetched from api
- [ ] Basic CI-able checks: `go vet`/`go test`, `tsc --noEmit`, `next build`

**Demo:** `make dev` → localhost shows a placeholder page that round-trips through the Go API.

## Phase 1 — Content model + public site in the new design

Goal: airfa.pt fully reproduced in the new design, content served from the DB (seeded), no CMS UI yet.

- [ ] Migrations: pages, page_revisions, menus, settings, media, events, posts, partners, activities
- [ ] Public API: `GET /v1/site`, `GET /v1/pages/:slug`, collections endpoints
- [ ] Media serving with image variants
- [ ] Design tokens in Tailwind (colors, fonts per design-spec §1 — sample exact values from mockups)
- [ ] Global chrome: utility bar, header + dropdown menus, red secondary bar, footer, back-to-top — all driven by menus/settings data
- [ ] Block registry + all Phase-1 blocks (design-spec §6): hero-slider, section-band, events-fan-carousel, activities-slider, blog-grid, partners-strip, page-hero-light, page-hero-image, story-era, gallery-masonry, documents-showcase, documents-table, rich-text, image-text, contact-info
- [ ] Homepage matches mockup-01; Banda matches mockup-02; Estatutos matches mockup-03 (pixel-close review with owner)
- [ ] Remaining pages built from blocks (História, Órgãos Sociais, Hino, Biblioteca, Documentos, Actividades, Salas de Espetáculo, Cine-Teatro, Sala de Cinema, Contactos, 404)
- [ ] Seed script: menus, settings, all pages with real airfa.pt content (content-inventory.md), initial media import
- [ ] Responsive pass (mobile/tablet) + accessibility pass
- [ ] SEO: metadata per page, sitemap.xml, robots.txt, Open Graph

**Demo:** the complete public site, new design, all real content — visually approved by owner.

## Phase 2 — CMS (Área Reservada)

Goal: owner can edit everything without a developer.

- [ ] Auth: login page, sessions, argon2id, rate limiting; seed admin user; middleware guard on `/area-reservada`
- [ ] Users management (Admin): create/disable users, roles Admin/Editor
- [ ] Admin shell: navigation, list of pages, site settings form
- [ ] Page editor: block list with drag-and-drop reorder, add/remove/duplicate blocks, per-block edit forms, internal-link picker
- [ ] Draft/publish workflow + preview of drafts
- [ ] Publish → version snapshot; version history UI + restore
- [ ] Menu editor (all zones, nesting, reorder)
- [ ] Media library: upload, browse, replace, delete, alt text; PDF metadata for documents blocks
- [ ] Collections editors: events, posts, partners, activities
- [ ] Page management: create/rename/delete (soft) pages, slugs, SEO fields
- [ ] Publish triggers public-site revalidation (webhook)
- [ ] Round-trip Playwright test: login → edit → publish → public page updated

**Demo:** owner edits a page, adds a block, publishes, sees it live; restores an old version.

## Phase 3 — Deployment & operations

- [ ] Dockerfiles (api static binary; web standalone) + compose file for Dokploy
- [ ] Dokploy setup on VPS: services, volumes (pgdata, media), domain + TLS, env secrets
- [ ] Staging-then-production cutover plan for www.airfa.pt DNS
- [ ] Automated backups: nightly pg_dump + media volume, ≥14-day retention; documented restore procedure (test it once)
- [ ] Uptime/basic monitoring (Dokploy health checks)
- [ ] Handover doc for the owner: how to log in, edit, publish, restore, and who to call

**Demo:** site live on www.airfa.pt, owner-editable, backed up.

## Later / nice-to-have (explicitly out of v1)

- Site search (header icon works on-page content)
- Multi-language (schema is ready for a locale dimension)
- Cookieless analytics
- Newsletter integration
- Image focal-point cropping in CMS
