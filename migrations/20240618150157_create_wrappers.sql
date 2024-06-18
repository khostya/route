-- +goose Up
-- +goose StatementBegin
create table if not exists ozon.wrappers
(
    order_id       text primary key references ozon.orders (id),
    type           ozon.wrapper not null,
    price_in_rub   numeric       not null,
    capacity_in_kg float4       not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists ozon.wrappers;
-- +goose StatementEnd
