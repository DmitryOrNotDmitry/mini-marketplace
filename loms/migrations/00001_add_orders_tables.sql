-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    order_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    status TEXT NOT NULL DEFAULT 'new'::text
);

CREATE TABLE order_items (
    id BIGSERIAL PRIMARY KEY,
    sku BIGINT NOT NULL,
    order_id BIGINT NOT NULL REFERENCES orders(order_id),
    count BIGINT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE order_items;
DROP TABLE orders;
-- +goose StatementEnd
