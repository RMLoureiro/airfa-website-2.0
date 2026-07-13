# Requirements

> **Living document.** When the owner answers an open question (see `STATUS.md`), replace the *(proposal — confirm)* marker with the confirmed decision and date. When scope changes, edit the requirement — don't let this file drift from reality.

## 1. Goal

Rebuild www.airfa.pt with the new design (mockups in `design/`) and give AIRFA a WordPress-like editing experience: a non-technical person logs in and can change essentially everything on the site — text, images, blocks, whole pages, menus — without a developer.

At launch the site must contain **all the information of the current airfa.pt** (see `content-inventory.md`), presented in the new design.

## 2. Public site

- Portuguese (pt-PT) only. Schema must not prevent adding locales later.
- Server-rendered (SEO, social link previews). Proper `<title>`/meta/Open Graph per page.
- Fully responsive (mockups are desktop; mobile layouts follow the same design language — collapsible nav, stacked sections).
- Pages are compositions of **typed blocks** (hero slider, rich text, gallery, documents table, partners strip, etc. — full catalog in `design-spec.md` §Blocks).
- Navigation (top utility bar, main menu with dropdowns, red secondary menu, footer link columns) is **data**, editable in the CMS — including nesting (tabs → sub-pages) and links to internal pages or external URLs.
- 404 page in the new design.
- Accessibility: semantic HTML, alt texts (editable per image in CMS), keyboard-navigable menus and carousels.

## 3. CMS — "Área Reservada"

Login-protected admin at `/area-reservada`. Everything below happens through the UI, no code.

*There are no mockups for the admin* — the strict-design rule applies to the public site only. The admin should be a clean, functional UI (shadcn/ui components, AIRFA red as accent), all labels in Portuguese, optimized for a non-technical editor: obvious buttons, confirmations before destructive actions, and no unexplained jargon.

### 3.1 Page management
- Create, rename, delete, and reorder pages; each page has a slug, SEO fields, and a block list.
- Add / remove / reorder / duplicate blocks on a page; edit each block's fields through forms (text, images, links, lists of items).
- Internal linking: when a field is a link, the editor can pick an existing page from a list (or paste an external URL).
- Draft vs published: edits are saved as draft and only visible publicly after "Publicar".
- Preview drafts before publishing.

### 3.2 Menus
- Edit all menu zones (main nav + dropdowns, secondary red bar, utility top bar, footer columns).
- Menu items point to internal pages or external URLs; support nesting for dropdowns; drag to reorder.

### 3.3 Media library
- Upload images/PDFs; browse, search, replace, delete.
- Images get automatic resized variants (thumbnail/medium/large) for performance.
- PDFs power the documents-table block (Estatutos-style page: name, size, modified date, download).

### 3.4 Content collections (structured content beyond pages)
- **Events / Cartazes** — poster image, title, date(s), optional description + link. Feed the homepage "Próximos eventos em cartaz" carousel and an events page.
- **Blog posts** — title, cover, excerpt, body (blocks), publish date. Feed the homepage Blog grid. *(Confirmed real content, 2026-07-12 — the association will maintain posts; build + seed fully in Phase 1.)*
- **Partners** — logo + name + optional URL. Feed the "Parceiros e Apoios" strip shown on every page.
- **Activities** — name, category, image, schedule/instructor info. Feed the activities slider and pages.

### 3.5 Site settings
- Contact info (address, email, phone, NIB), social links, DOAR button target (v1: an on-site how-to-donate / IBAN page — internal page link, not an external platform), footer texts, logo/favicon.

## 4. Users & roles *(proposal — confirm)*

- Named accounts with email + password. No self-registration; admins create accounts.
- **Admin**: everything, including user management and site settings.
- **Editor**: edit content (pages, blocks, menus, media, collections) and publish; cannot manage users or delete version history.
- Passwords hashed (argon2id); sessions via httpOnly cookies; rate-limited login.

## 5. Versioning & safety *(proposal — confirm)*

- Every publish creates a **version snapshot** of the affected page (and of menus/settings when those change).
- Version history per page with "restore this version" (restore creates a new draft; publishing it makes it live).
- Retention: keep last N versions per page (e.g. 20) — simple, predictable.
- Deleting a page requires confirmation and is soft-delete (recoverable by Admin for 30 days).

## 6. Non-functional

- Fast: public pages cached/ISR'd; CMS changes visible on the public site within seconds of publishing.
- Backups: automated Postgres dumps + media volume backup (Phase 3, Dokploy).
- No analytics/cookies beyond what's strictly needed (association site; keep GDPR surface minimal). If analytics are added later, prefer a cookieless option.
- The system must run on a small VPS (single node, Docker).

## 7. Explicit non-goals (v1)

- No multi-language UI/content (schema-ready only).
- No e-commerce/payments (DOAR is a link, not a checkout).
- No member portal/user accounts for the public.
- No comments on blog posts.
