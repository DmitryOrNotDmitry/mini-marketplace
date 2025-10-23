-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders_event_outbox (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT,
    order_status TEXT NOT NULL,
    moment TIMESTAMP NOT NULL,
    event_status TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders_event_outbox;
-- +goose StatementEnd
