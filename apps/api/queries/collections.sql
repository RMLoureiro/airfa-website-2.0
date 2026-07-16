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