# Content Inventory — what must exist at launch

Source of truth for current content: the live site **https://www.airfa.pt** and its codebase at **`../airfa-website/`** (static HTML + assets; content is partly rendered by its own mini-CMS, so extract text from the live pages or the HTML files there). Use it for **content only** — not design, not architecture.

## Site map of the current site (to be reproduced, reorganized per new design)

| Current page | File in `../airfa-website/` | Notes for 2.0 |
|---|---|---|
| Início | `index.html` | Becomes new homepage (mockup-01) |
| A Academia (dropdown parent) | `a-academia.html` | Section landing under "Quem Somos" |
| — História | `historia.html` | Story-style page (mockup-02 pattern fits: timeline/era blocks) |
| — Estatutos | `estatutos.html` | Documents page — mockup-03 is literally this page |
| — Órgãos Sociais | `orgaos-sociais.html` | Rich text / people lists |
| — Hino | `hino.html` | Lyrics + possibly audio |
| — Biblioteca | `biblioteca.html` | Rich text + images |
| — Documentos | `documentos.html` | Documents-table block (same as Estatutos) |
| Actividades | `actividades.html` | Activities collection + slider (mockup-01 "Atividades 2025/26") |
| Banda | `banda.html` | Mockup-02 is this page (legado de 1895, galeria) |
| Salas de Espetáculo (dropdown parent) | `salas-de-espetaculo.html` | |
| — Cine-Teatro | `cine-teatro.html` | |
| — Sala de Cinema | `sala-de-cinema.html` | |
| Contactos | `contactos.html` | Address, email, phone, NIB, map |
| 404 | `404.html` | New-design 404 |

## Navigation in the new design (from mockups)

- **Utility top bar**: Blog · Ajuda & Suporte · Imprensa · Login
- **Main nav**: Quem Somos ▾ · Actividades ▾ (labelled "O que fazemos" on one mockup — confirm final label) · Banda ▾ · Salas de Espetáculo ▾ · Contactos ▾ · search icon
- **Secondary red bar**: Missão/Atividades · Cartazes de Eventos · A nossa História · Junte-se a nós · **DOAR** (pill button)
- **Footer columns**: Explorar (Home, Sobre nós, Menu, Contactos) · Serviços (Encomendas, Reservas, Personalização) · Social (Instagram, Facebook, Twitter/X, LinkedIn) — footer content in mockups looks partly placeholder; owner to confirm real links.

Mapping old pages into these menus is an editorial decision for the owner during Phase 1 seeding; menus are CMS-editable anyway.

## Fixed facts (from mockups / current site — verify at seeding time)

- Name: **Academia de Instrução e Recreio Familiar Almadense (A.I.R.F.A.)**, "Desde 1895".
- Address: Rua Capitão Leitão, nº64, 2800-068 Almada.
- E-mail: academia.almadense@gmail.com
- Telefone: 21 272 9750 / 55
- NIB: 0036 0229 9910 0147 2274 6
- Banda founded 1895; maestro since 2011: Francisco Pinto (per mockup-02 text — verify against live site).
- Estatutos PDFs: 1927, 1967, 1987 (mockup-03 table).

## Assets to migrate

- Crest/logo: `agent/design/logo-airfa.png` (clean master). Old site favicons/crest in `../airfa-website/assets/`.
- Photo galleries: `../airfa-website/assets/galeria/` and related folders — review and import the good ones into the media library during seeding.
- PDFs (estatutos, documentos): collect from live site / old repo at seeding time.
