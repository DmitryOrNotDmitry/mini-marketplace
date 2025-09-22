-- +goose Up
-- +goose StatementBegin
INSERT INTO stocks (sku, total_count, reserved) VALUES
(139275865, 65534, 0),
(2956315, 300, 30),
(1076963, 300, 35),
(135717466, 100, 20),
(135937324, 300, 30),
(1625903, 10000, 0),
(1148162, 300, 0),
(139819069, 100, 100),
(2618151, 300, 0),
(2958025, 300, 0),
(3596599, 300, 0),
(3618852, 300, 0),
(4288068, 300, 0),
(4465995, 300, 0),
(30816475, 300, 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM stocks;
-- +goose StatementEnd
