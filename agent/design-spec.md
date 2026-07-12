# Design Specification

**Source of truth:** the mockups in `design/`. This document describes them so agents can plan and name things, but when in doubt, open the images. Follow the design **strictly**.

- `design/mockup-01-homepage.jpeg` — Homepage
- `design/mockup-02-banda.png` — Banda (content/story page pattern)
- `design/mockup-03-estatutos.png` — Estatutos (documents page pattern)
- `design/logo-airfa.png` — A.I.R.F.A. crest (red shield, white diagonal band with "A.I.R.F.A.", gold mural crown)

## 1. Design tokens (extract exact values from mockups during Phase 1)

**Colors** (approximate — sample the images for finals):
- `red-primary` — AIRFA red (crest / red bars / copyright bar), ≈ `#D5232A`
- `black-ink` — near-black for dark sections & headings, ≈ `#0B0B0C`
- `white` — backgrounds, text on dark
- `gold` — crest crown gold, ≈ `#E8A33D`
- `orange-accent` — GALERIA band / activity slide background, ≈ `#F5A83C`
- `footer-blue` — pale desaturated blue-grey footer background, ≈ `#C9D4DF` (with a large watermark of the crest, low-opacity, bottom-left)
- `grey-placeholder` — light grey used in card/image placeholders, ≈ `#D8DCE0`

**Typography** (mockups appear to use a bold grotesque for headings; confirm exact family with owner):
- Headings/UI: a tight bold sans (proposal: **Archivo** or **Inter**, weights 600–900). Section titles are UPPERCASE, very bold, large (e.g. "ATIVIDADES 2025/26", "BLOG", "PARCEIROS E APOIOS" — left-aligned; "PRÓXIMOS EVENTOS EM CARTAZ", "GALERIA" — centered white on colored band).
- Display numerals/hero display: the giant "1895" on mockup-02 and huge "BANDA" title — serif display for numerals (proposal: **Playfair Display** or similar Didone) and ultra-bold sans for page titles.
- Body: same sans, 400/500, comfortable line-height.

**Shape language:** generous rounded corners on cards/images (~16–24px); pill-shaped buttons and tags; full-bleed colored bands for section separators.

## 2. Global chrome (all pages)

### 2.1 Utility top bar
Very slim, light grey. Right-aligned small uppercase links: `BLOG · AJUDA & SUPORTE · IMPRENSA · LOGIN`. (CMS: utility menu zone.)

### 2.2 Main header (white)
- Left: crest logo + wordmark "AIRFA" with tiny "Desde 1895" underneath.
- Center/right: uppercase menu items with dropdown carets: `QUEM SOMOS ▾ · ACTIVIDADES ▾ (a.k.a. "O QUE FAZEMOS" on one mockup — final label TBD) · BANDA ▾ · SALAS DE ESPETÁCULO ▾ · CONTACTOS ▾` + search icon (opens search overlay; search can be Phase 2+).
- Dropdowns host sub-pages (e.g. Quem Somos → História, Estatutos, Órgãos Sociais, Hino, Biblioteca, Documentos).

### 2.3 Secondary red bar
Full-width `red-primary` bar directly under the header. White links: contextual quick links (`Missão / Atividades · Cartazes de Eventos · A nossa História · Junte-se a nós`) and on the right a **white pill button "♡ DOAR"** (red text/icon). (CMS: secondary menu zone + DOAR target in settings.)

### 2.4 Footer
- Big `footer-blue` area, large low-opacity crest watermark on the left; the actual crest logo displayed prominently left of the link columns.
- Link columns with small uppercase grey headers: `EXPLORAR`, `SERVIÇOS`, `SOCIAL` (mockup shows duplicated columns — treat as N configurable columns).
- `ONDE ESTAMOS` block: association full name, address, e-mail, telephone, NIB; a square placeholder beside it (map or photo — TBD).
- Circular back-to-top button (black, chevron up) bottom-right.
- Bottom bar in `red-primary`: white text `© 2026 ACADEMIA ALMADENSE DE INSTRUÇÃO E RECRIAÇÃO FAMILIAR - TODOS OS DIREITOS RESERVADOS.` left, `DESENVOLVIDO POR Square²` right.

## 3. Homepage (mockup-01)

Top to bottom:

1. **Hero slider (dark)** — full-width black section with a background photo right (theatre stage). Left: big white headline with the last word underlined ("Um espetáculo mais do que <u>inesquecível</u>"), two pill CTAs side by side (outline "SABER MAIS" + solid white "SABER MAIS" — i.e., primary/secondary variants, editable labels/links). Below: a **news strip** of 4 compact items (title + 2-line grey excerpt), the active item marked with a red progress bar on top; a pause/play control on the left. Items rotate with the hero. (CMS block: `hero-slider` with N slides: image, headline, highlighted word, CTAs, linked news items.)
2. **Red band title** — full-width `red-primary` band, centered bold white uppercase title: "PRÓXIMOS EVENTOS EM CARTAZ". (`section-band`, editable text/color.)
3. **Events fan carousel** — event/modalidade posters fanned like held playing cards (center cards face-on, side cards rotated), floating name tags above ("JAM IN DA HOUSE" pink pill, "XADREZ" black pill). Round prev/next arrows at the sides. Fed by the Events collection. (`events-fan-carousel`.)
4. **Activities slider** — left-aligned bold title "ATIVIDADES 2025/26", then a full-width rounded card slider; each slide is a colored panel (e.g. orange) with a huge title ("DESPORTO"), a white pill CTA ("VER TODAS AS ATIVIDADES") and a person photo on the right. Dots bottom-left, round arrows bottom-right. Fed by Activities categories. (`activities-slider`.)
5. **Blog grid** — bold title "BLOG", three large rounded cards (image + eventually title/excerpt; mockup shows placeholders). Fed by latest posts. (`blog-grid`.)
6. **Partners strip** — bold title "PARCEIROS E APOIOS", a horizontal row of pill-shaped logo placeholders (marquee/scroll if overflowing). Fed by Partners collection. (`partners-strip` — appears on virtually every page above the footer.)

## 4. Content/story page — Banda (mockup-02)

Pattern for rich storytelling pages (Banda, História, …):

1. **Light page hero** — huge black uppercase title ("BANDA") centered on light background, small grey caps subtitle underneath ("LEGADO DE 1895"). (`page-hero-light`.)
2. **Story/era section (dark)** — full-bleed black section with a blurred historic photo background: left-aligned white heading ("A banda XPTO NASCE"), intro paragraph, then a **giant serif year numeral "1895"** spanning the width, followed by more body text. On the right edge: a vertical **timeline scrubber** (tick marks + current year label) suggesting multiple eras the user can scroll/jump through; a small red icon button top-right (view image). (`story-era` block: N eras each with year, heading, text, background image.)
3. **Orange band title** — `orange-accent` full-width band, centered white "GALERIA". (`section-band` with color option.)
4. **Gallery masonry** — rounded-corner image grid in a masonry arrangement (mixed tile sizes/heights). Fed by a media gallery. (`gallery-masonry`.)
5. Partners strip + footer as global.

## 5. Documents page — Estatutos (mockup-03)

Pattern for document-download pages (Estatutos, Documentos):

1. **Dark image hero** — full-width dark photo banner (theatre seats) with left-aligned white title "Estatutos". (`page-hero-image`.)
2. **Fanned document preview** — the documents displayed as slightly rotated overlapping sheets with their years ("1927", "1967", "1987") handwritten-style above each. (`documents-showcase`; optional visual companion to the table.)
3. **Documents table** — columns `NOME · TAMANHO · DATA MODIFICAÇÃO · AÇÕES`; grey uppercase header row; bold black file names; blue-grey size/date text; download icon button per row. Fed by PDFs in the media library. (`documents-table`.)
4. Partners strip + footer as global.

## 6. Block catalog (initial)

Global chrome (header/menus/footer) is site-level, not blocks. Page-level blocks:

| Block type | Used in | Key fields |
|---|---|---|
| `hero-slider` | Home | slides[]: bg image, headline, underlinedWord, ctas[], newsItems[] (title, excerpt, link) |
| `section-band` | Home, Banda | text, background color (red/orange/black), alignment |
| `events-fan-carousel` | Home | source: events collection (auto) or manual picks |
| `activities-slider` | Home | slides[]: title, color, image, cta |
| `blog-grid` | Home | count, source: latest posts |
| `partners-strip` | all pages | source: partners collection |
| `page-hero-light` | Banda-style | title, subtitle |
| `page-hero-image` | Estatutos-style | title, bg image, overlay strength |
| `story-era` | Banda, História | eras[]: year, heading, body, bg image |
| `gallery-masonry` | Banda | images[] (media refs, alt) |
| `documents-showcase` | Estatutos | docs[] (media refs, year label) |
| `documents-table` | Estatutos, Documentos | docs[] (media refs; size/date auto) |
| `rich-text` | any | markdown/portable text |
| `image` / `image-text` | any | media ref, alt, caption, layout |
| `contact-info` | Contactos | fields from settings + free text, map embed |

Add new block types as pages demand them (Hino, Órgãos Sociais, Salas…) — extend this table when you do.

## 7. Interaction notes

- Carousels: auto-advance hero (with pause button per mockup), manual arrows elsewhere; all swipeable on touch; keyboard accessible.
- Hover states: menu items underline/darken; cards lift slightly; buttons invert (white↔red).
- Back-to-top button appears after scrolling past the first section.
- Mobile: header collapses to hamburger (utility links inside the drawer); red bar becomes horizontally scrollable chips or lives in the drawer; fan carousel simplifies to a swipe carousel; tables become stacked cards or horizontally scrollable.
