-- +goose Up
-- +goose StatementBegin
create schema if not exists ozon;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop schema if exists ozon;
-- +goose StatementEnd
