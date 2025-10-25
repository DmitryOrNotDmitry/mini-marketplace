-- name: AddComment :one
insert into comments(id, user_id, sku, content, created_at)
values (nextval('comment_id_manual_seq') + $1, $2, $3, $4, $5)
returning id;

-- name: GetCommentByID :one
select *
from comments
where id = $1;

-- name: GetCommentByIDForUpdate :one
select *
from comments
where id = $1 for update;

-- name: UpdateContent :exec
update comments
set content = $2
where id = $1;

-- name: GetCommentsBySKU :many
select *
from comments
where sku = $1;

-- name: GetCommentsByUser :many
select *
from comments
where user_id = $1;
