-- name: GetPublishedPageBySlug :one
select p.slug, p.title, p.seo, r.blocks
from pages p
join page_revisions r on r.page_id = p.id and r.kind = 'published'
where p.slug = $1 and p.deleted_at is null;