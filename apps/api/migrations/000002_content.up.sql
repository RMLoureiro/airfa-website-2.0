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