-- +goose Up
-- +goose StatementBegin
create index orders_created_at_idx on ozon.orders using btree(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists ozon.orders_created_at_idx;
-- +goose StatementEnd
