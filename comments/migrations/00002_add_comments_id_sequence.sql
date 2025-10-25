-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE comment_id_manual_seq INCREMENT 1000 START 1000;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SEQUENCE comment_id_manual_seq;
-- +goose StatementEnd
