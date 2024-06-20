-- +goose NO TRANSACTION
-- +goose Up
-- +goose StatementBegin
create index concurrently if not exists orders_recipient_id_idx on ozon.orders using btree(recipient_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists ozon.orders_recipient_id_idx;
-- +goose StatementEnd
