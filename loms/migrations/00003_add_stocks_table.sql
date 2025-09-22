-- +goose Up
-- +goose StatementBegin
CREATE TABLE stocks (
    sku BIGINT PRIMARY KEY,
    total_count BIGINT NOT NULL,
    reserved BIGINT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stocks;
-- +goose StatementEnd
