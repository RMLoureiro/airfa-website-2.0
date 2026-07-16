-- name: GetSettings :one
select data from settings limit 1;

-- name: ListMenus :many
select zone, items from menus;