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