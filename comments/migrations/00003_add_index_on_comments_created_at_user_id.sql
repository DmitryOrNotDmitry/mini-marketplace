-- +goose Up
-- +goose NO TRANSACTION
CREATE INDEX CONCURRENTLY idx_comments_sku_created_at_desc_user_id_asc
ON comments(sku, created_at DESC, user_id ASC);

-- +goose Down
-- +goose NO TRANSACTION
DROP INDEX CONCURRENTLY idx_comments_sku_created_at_desc_user_id_asc;
