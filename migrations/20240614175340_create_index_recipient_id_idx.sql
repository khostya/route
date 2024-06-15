-- +goose Up
-- +goose StatementBegin
create index orders_recipient_id_idx on ozon.orders using btree(recipient_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists ozon.orders_recipient_id_idx;
-- +goose StatementEnd
